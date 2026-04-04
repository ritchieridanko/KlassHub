package usecases

import (
	"context"
	"errors"
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
	GetSession(ctx context.Context, refreshToken string) (s *models.Session, err *ce.Error)
	RefreshSession(ctx context.Context, req *models.RefreshSessionReq) (at *models.AuthToken, err *ce.Error)
	RevokeSession(ctx context.Context, refreshToken string) (err *ce.Error)
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

func NewSessionUsecase(appName string, accessTokenExpiry, refreshTokenExpiry time.Duration, sr repositories.SessionRepository, tx *database.Transactor, v *validator.Validator, j *jwt.JWT) SessionUsecase {
	return &sessionUsecase{
		appName:            appName,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
		sr:                 sr,
		transactor:         tx,
		validator:          v,
		jwt:                j,
	}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, req *models.CreateSessionReq) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.CreateSession")
	defer span.End()

	authIDField := logger.NewField("auth_id", req.AuthID)
	schoolIDField := logger.NewField("school_id", req.SchoolID)
	roleField := logger.NewField("role", req.Role)

	// UUID Creation
	uuid := utils.GenerateUUID()

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
		if transportCtx == nil {
			return ce.NewError(
				ce.CodeMissingContextValue,
				ce.MsgInternalServer,
				errors.New("transport missing from context"),
				authIDField,
				schoolIDField,
				roleField,
			)
		}

		// Active Session Revocation
		sessionID, err := u.sr.RevokeActive(
			ctx,
			&models.RevokeActiveSessionParams{
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
			txErr.Unwrap(),
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

func (u *sessionUsecase) GetSession(ctx context.Context, refreshToken string) (*models.Session, *ce.Error) {
	return u.sr.GetByRefreshToken(ctx, refreshToken)
}

func (u *sessionUsecase) RefreshSession(ctx context.Context, req *models.RefreshSessionReq) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.RefreshSession")
	defer span.End()

	authIDField := logger.NewField("auth_id", req.AuthID)
	schoolIDField := logger.NewField("school_id", req.SchoolID)
	roleField := logger.NewField("role", req.Role)

	// UUID Creation
	uuid := utils.GenerateUUID()

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
		// Session Revocation
		s, err := u.sr.Revoke(
			ctx,
			&models.RevokeSessionParams{
				RefreshToken: req.RefreshToken,
				ExpiresAt:    now,
			},
		)
		if err != nil {
			return err.Append(authIDField, schoolIDField, roleField)
		}
		if s.AuthID != req.AuthID {
			return ce.NewError(
				ce.CodeSessionNotOwned,
				ce.MsgInvalidSession,
				nil,
				authIDField,
				logger.NewField("session_auth_id", s.AuthID),
				schoolIDField,
				roleField,
			)
		}

		// Session Creation
		err = u.sr.Create(
			ctx,
			&models.CreateSessionData{
				ParentID:     &s.ID,
				AuthID:       s.AuthID,
				RefreshToken: uuid.String(),
				UserAgent:    s.UserAgent,
				IPAddress:    s.IPAddress,
				ExpiresAt:    now.Add(u.refreshTokenExpiry),
			},
		)
		if err != nil {
			return err.Append(schoolIDField, roleField)
		}
		return nil
	})
	if txErr != nil && txErr.Code() == ce.CodeDBTransaction {
		return nil, ce.NewError(
			txErr.Code(),
			txErr.Message(),
			txErr.Unwrap(),
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

func (u *sessionUsecase) RevokeSession(ctx context.Context, refreshToken string) *ce.Error {
	_, err := u.sr.Revoke(
		ctx,
		&models.RevokeSessionParams{
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().UTC(),
		},
	)
	return err
}
