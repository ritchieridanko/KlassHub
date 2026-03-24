package usecases

import (
	"context"
	"fmt"
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
	CreateSession(ctx context.Context, req *models.CreateSessionReq) (at *models.AuthToken, err *ce.Error)
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

func (u *sessionUsecase) CreateSession(ctx context.Context, req *models.CreateSessionReq) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.CreateSession")
	defer span.End()

	authIDField := logger.NewField("auth_id", req.AuthID)
	schoolIDField := logger.NewField("school_id", req.SchoolID)
	roleField := logger.NewField("role", req.Role)

	// UUID Creation
	uuid, err := utils.GenerateUUIDv7()
	if err != nil {
		return nil, ce.NewError(
			ce.CodeUUIDGenerationFailed,
			ce.MsgInternalServer,
			err,
			authIDField,
			schoolIDField,
			roleField,
		)
	}

	// JWT Creation
	now := time.Now().UTC()
	jwt, err := u.jwt.Generate(req.AuthID, req.SchoolID, req.Role, req.IsVerified, &now)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeJWTGenerationFailed,
			ce.MsgInternalServer,
			err,
			authIDField,
			schoolIDField,
			roleField,
		)
	}

	txErr := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		transportCtx := utils.CtxTransport(ctx)

		// Active Session Revocation
		sessionID, err := u.sr.RevokeActive(
			ctx,
			&models.RevokeSessionParams{
				AuthID:    req.AuthID,
				UserAgent: transportCtx.UserAgent,
				IPAddress: transportCtx.IPAddress,
				ExpiresAt: now,
			},
		)
		if err != nil {
			return err.Append(schoolIDField, roleField)
		}

		// Session Creation
		// NOTE: Set revoked session (if any) as parent session to maintain session lineage
		data := models.CreateSessionData{
			AuthID:       req.AuthID,
			RefreshToken: uuid.String(),
			UserAgent:    transportCtx.UserAgent,
			IPAddress:    transportCtx.IPAddress,
			ExpiresAt:    now.Add(u.refreshTokenExpiry),
		}
		if invalidSessionID := int64(0); sessionID > invalidSessionID {
			data.ParentID = &sessionID
		}
		if err := u.sr.Create(ctx, &data); err != nil {
			return err.Append(schoolIDField, roleField)
		}
		return nil
	})
	if txErr != nil && txErr.Code() == ce.CodeDBTransaction {
		return nil, ce.NewError(
			txErr.Code(),
			txErr.Message(),
			fmt.Errorf("failed to create session: %w", txErr.Unwrap()),
			authIDField,
			schoolIDField,
			roleField,
		)
	}
	if txErr != nil {
		return nil, txErr
	}

	return &models.AuthToken{
		AccessToken: &models.AccessToken{
			Token:     jwt,
			ExpiresIn: int64(u.accessTokenExpiry.Seconds()),
		},
		RefreshToken: &models.RefreshToken{
			Token:     uuid.String(),
			ExpiresIn: int64(u.refreshTokenExpiry.Seconds()),
		},
	}, nil
}
