package usecases

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
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
	Login(ctx context.Context, req *models.LoginReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	CreateSchoolAuth(ctx context.Context, req *models.CreateSchoolAuthReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
}

type authUsecase struct {
	appName    string
	su         SessionUsecase
	ar         repositories.AuthRepository
	transactor *database.Transactor
	validator  *validator.Validator
	bcrypt     *bcrypt.BCrypt
	logger     *logger.Logger
}

func NewAuthUsecase(appName string, su SessionUsecase, ar repositories.AuthRepository, tx *database.Transactor, v *validator.Validator, b *bcrypt.BCrypt, l *logger.Logger) AuthUsecase {
	return &authUsecase{
		appName:    appName,
		su:         su,
		ar:         ar,
		transactor: tx,
		validator:  v,
		bcrypt:     b,
		logger:     l,
	}
}

func (u *authUsecase) Login(ctx context.Context, req *models.LoginReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.Login")
	defer span.End()

	// Data Normalization
	identifier := utils.NormalizeString(req.Identifier)

	// Data Validation
	if ok, why := u.validator.Identifier(identifier); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Auth Fetching
	a, err := u.ar.GetByIdentifier(ctx, identifier)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return nil, nil, ce.NewError(
			ce.CodeIdentifierNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			err.Fields()...,
		)
	}
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)
	schoolIDField := logger.NewField("school_id", a.SchoolID)
	roleField := logger.NewField("role", a.Role)

	// Password Validation
	if err := u.bcrypt.Validate(a.Password, req.Password); err != nil {
		return nil, nil, ce.NewError(
			ce.CodeWrongPassword,
			ce.MsgInvalidCredentials,
			err,
			authIDField,
			schoolIDField,
		)
	}

	// Role and Subdomain Validation
	// NOTE:
	// - Students and Instructors can only login from LMS
	// - Administrators and Schools can only login from Admin
	if subdomain := utils.CtxSubdomain(ctx); !u.validator.RoleAllowedSubdomain(a.Role, subdomain) {
		return nil, nil, ce.NewError(
			ce.CodeWrongSubdomain,
			ce.MsgInvalidCredentials,
			ce.ErrWrongSubdomain,
			authIDField,
			schoolIDField,
			roleField,
			logger.NewField("subdomain", subdomain),
		)
	}

	// Session Creation
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionReq{
			AuthID:     a.ID,
			SchoolID:   a.SchoolID,
			Role:       a.Role,
			IsVerified: a.IsVerified(),
		},
	)

	return a, at, err
}

func (u *authUsecase) CreateSchoolAuth(ctx context.Context, req *models.CreateSchoolAuthReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.CreateSchoolAuth")
	defer span.End()

	// Data Normalization
	email := utils.NormalizeString(req.Email)

	// Data Validation
	if ok, why := u.validator.Email(email); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}
	if ok, why := u.validator.Password(req.Password); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	var a *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Email Availability Check
		available, err := u.ar.IsEmailAvailable(ctx, email)
		if err != nil {
			return err
		}
		if !available {
			return ce.NewError(ce.CodeEmailNotAvailable, ce.MsgEmailNotAvailable, nil)
		}

		// Password Hashing
		hash, hashErr := u.bcrypt.Hash(req.Password)
		if hashErr != nil {
			return ce.NewError(ce.CodeBCryptHashingFailed, ce.MsgInternalServer, hashErr)
		}

		// Auth Creation
		schoolProfileNotYetExists := int64(0)
		a, err = u.ar.Create(
			ctx,
			&models.CreateAuthData{
				SchoolID: schoolProfileNotYetExists,
				Email:    &email,
				Password: hash,
				Role:     constants.RoleSchool,
			},
		)
		return err
	})
	if err != nil && err.Code() == ce.CodeDBTransaction {
		return nil, nil, ce.NewError(
			err.Code(),
			err.Message(),
			fmt.Errorf("failed to create school auth: %w", err.Unwrap()),
		)
	}

	// Session Creation
	// NOTE: Fail to create session does not fail create school auth process
	at, err := u.su.CreateSession(
		ctx,
		&models.CreateSessionReq{
			AuthID:     a.ID,
			SchoolID:   a.SchoolID,
			Role:       a.Role,
			IsVerified: a.IsVerified(),
		},
	)
	if err != nil {
		u.logger.Warn(
			ctx,
			"created school auth. failed to create session",
			err.Append(
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err),
			).Fields()...,
		)
		return a, nil, nil
	}

	return a, at, nil
}
