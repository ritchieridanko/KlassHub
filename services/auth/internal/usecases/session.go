package usecases

import (
	"context"
	"time"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/jwt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/validator"
	"go.opentelemetry.io/otel"
)

type SessionUsecase interface {
	CreateSession(ctx context.Context, req *models.CreateSessionRequest) (at *models.AuthToken, err *ce.Error)
}

type sessionUsecase struct {
	appName            string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	sr                 repositories.SessionRepository
	transactor         *database.Transactor
	validator          *validator.Validator
	jwt                *jwt.JWT
}

func NewSessionUsecase(appName string, accessTokenExpiry, refreshTokenExpiry time.Duration, sr repositories.SessionRepository, tx *database.Transactor, v *validator.Validator, jwt *jwt.JWT) SessionUsecase {
	return &sessionUsecase{
		appName:            appName,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		sr:                 sr,
		transactor:         tx,
		validator:          v,
		jwt:                jwt,
	}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, req *models.CreateSessionRequest) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.CreateSession")
	defer span.End()

	authIDField := logger.NewField("auth_id", req.AuthID)
	schoolIDField := logger.NewField("school_id", req.SchoolID)

	// Data Normalization
	ua := utils.NormalizeString(req.RequestMeta.UserAgent)
	ip := utils.NormalizeString(req.RequestMeta.IPAddress)

	// Data Validation
	if ok, why := u.validator.UserAgent(ua); !ok {
		return nil, ce.NewError(
			ce.CodeInvalidRequestMeta,
			why,
			nil,
			authIDField,
			schoolIDField,
		)
	}
	if ok, why := u.validator.IPAddress(ip); !ok {
		return nil, ce.NewError(
			ce.CodeInvalidRequestMeta,
			why,
			nil,
			authIDField,
			schoolIDField,
		)
	}

	// UUID Creation
	uuid, err := utils.GenerateUUIDv7()
	if err != nil {
		return nil, ce.NewError(
			ce.CodeUUIDGenerationFailed,
			ce.MsgInternalServer,
			err,
			authIDField,
			schoolIDField,
		)
	}

	// JWT Creation
	now := time.Now().UTC()
	jwt, err := u.jwt.Generate(req.AuthID, req.SchoolID, req.Role, req.IsEmailVerified, &now)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeJWTGenerationFailed,
			ce.MsgInternalServer,
			err,
			authIDField,
			schoolIDField,
		)
	}

	te := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Active Session Revocation
		sessionID, err := u.sr.RevokeActive(
			ctx,
			&models.RevokeSession{
				AuthID:    req.AuthID,
				UserAgent: ua,
				IPAddress: ip,
			},
		)
		if err != nil {
			return err.AppendFields(schoolIDField)
		}

		/* Session Creation
		 * Note: Set revoked session (if any) as parent session
		 */
		data := models.CreateSession{
			AuthID:       req.AuthID,
			RefreshToken: uuid.String(),
			UserAgent:    ua,
			IPAddress:    ip,
			ExpiresAt:    now.Add(u.refreshTokenExpiry),
		}
		if invalidSessionID := int64(0); sessionID > invalidSessionID {
			data.ParentID = &sessionID
		}
		if err := u.sr.Create(ctx, &data); err != nil {
			return err.AppendFields(schoolIDField)
		}
		return nil
	})

	return &models.AuthToken{
		AccessToken:           jwt,
		RefreshToken:          uuid.String(),
		AccessTokenExpiresIn:  int64(u.accessTokenExpiry.Seconds()),
		RefreshTokenExpiresIn: int64(u.refreshTokenExpiry.Seconds()),
	}, te
}
