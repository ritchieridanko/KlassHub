package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/bcrypt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/validator"
	"go.opentelemetry.io/otel"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *models.LoginRequest) (a *models.Auth, at *models.AuthToken, err *ce.Error)
}

type authUsecase struct {
	appName    string
	su         SessionUsecase
	ar         repositories.AuthRepository
	transactor *database.Transactor
	validator  *validator.Validator
	bcrypt     *bcrypt.BCrypt
}

func NewAuthUsecase(appName string, su SessionUsecase, ar repositories.AuthRepository, tx *database.Transactor, v *validator.Validator, bcrypt *bcrypt.BCrypt) AuthUsecase {
	return &authUsecase{
		appName:    appName,
		su:         su,
		ar:         ar,
		transactor: tx,
		validator:  v,
		bcrypt:     bcrypt,
	}
}

func (u *authUsecase) Login(ctx context.Context, req *models.LoginRequest) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.Login")
	defer span.End()

	subdomainField := logger.NewField("subdomain", req.Subdomain)

	// Data Normalization
	identifier := utils.NormalizeString(req.Identifier)

	// Data Validation
	if ok, why := u.validator.Identifier(identifier); !ok {
		return nil, nil, ce.NewError(
			ce.CodeInvalidIdentifier,
			why,
			nil,
			subdomainField,
		)
	}

	// Auth Fetching
	a, err := u.ar.GetByIdentifier(ctx, identifier)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return nil, nil, ce.NewError(
			ce.CodeIdentifierNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			subdomainField,
		)
	}
	if err != nil {
		return nil, nil, ce.NewError(
			err.Code(),
			err.Message(),
			err.Unwrap(),
			subdomainField,
		)
	}

	// Password Validation
	if err := u.bcrypt.Validate(a.Password, req.Password); err != nil {
		return nil, nil, ce.NewError(
			ce.CodeWrongPassword,
			ce.MsgInvalidCredentials,
			err,
			logger.NewField("auth_id", a.ID),
			subdomainField,
		)
	}

	// Session Creation
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionRequest{
			AuthID:          a.ID,
			SchoolID:        a.SchoolID,
			Role:            a.Role,
			IsEmailVerified: a.IsEmailVerified(),
			RequestMeta:     req.RequestMeta,
		},
	)
	if err != nil {
		return nil, nil, err.AppendFields(subdomainField)
	}

	return a, at, nil
}
