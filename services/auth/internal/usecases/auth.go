package usecases

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ritchieridanko/klasshub/services/auth/internal/constants"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/bcrypt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/validator"
	"github.com/ritchieridanko/klasshub/shared/contract/events/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *models.LoginReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	Logout(ctx context.Context, refreshToken string) (err *ce.Error)
	CreateSchoolAuth(ctx context.Context, req *models.CreateSchoolAuthReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	ResendVerification(ctx context.Context) (email string, err *ce.Error)
	VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authUsecase struct {
	appName                 string
	verificationTokenExpiry time.Duration
	su                      SessionUsecase
	ar                      repositories.AuthRepository
	tr                      repositories.TokenRepository
	transactor              *database.Transactor
	acp                     *publisher.Publisher
	avrp                    *publisher.Publisher
	validator               *validator.Validator
	bcrypt                  *bcrypt.BCrypt
	logger                  *logger.Logger
}

func NewAuthUsecase(appName string, verificationTokenExpiry time.Duration, su SessionUsecase, ar repositories.AuthRepository, tr repositories.TokenRepository, tx *database.Transactor, acp, avrp *publisher.Publisher, v *validator.Validator, b *bcrypt.BCrypt, l *logger.Logger) AuthUsecase {
	return &authUsecase{
		appName:                 appName,
		verificationTokenExpiry: verificationTokenExpiry,
		su:                      su,
		ar:                      ar,
		tr:                      tr,
		transactor:              tx,
		acp:                     acp,
		avrp:                    avrp,
		validator:               v,
		bcrypt:                  b,
		logger:                  l,
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

func (u *authUsecase) Logout(ctx context.Context, refreshToken string) *ce.Error {
	// Session Revocation
	// NOTE: Invalid session does not fail logout process
	err := u.su.RevokeSession(ctx, strings.TrimSpace(refreshToken))
	if err != nil && err.Code() != ce.CodeSessionNotFound {
		return err
	}
	return nil
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
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)
	roleField := logger.NewField("role", a.Role)

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
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
	}

	// Verification Token Creation
	// NOTE: Fail to create verification token does not fail create school auth process
	token := utils.GenerateUUID()
	err = u.tr.CreateVerification(
		ctx,
		&models.CreateVerificationTokenData{
			AuthID:   a.ID,
			Token:    token.String(),
			Duration: u.verificationTokenExpiry,
		},
	)
	if err != nil {
		u.logger.Warn(
			ctx,
			"created school auth. failed to create verification token",
			err.Append(
				roleField,
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
		return a, at, nil
	}

	eventTopicField := logger.NewField("event_topic", constants.EventTopicAC)

	// Event Publishing
	// NOTE: Fail to publish event does not fail create school auth process
	pubErr := u.acp.Publish(
		ctx,
		"auth_"+strconv.FormatInt(a.ID, 10),
		&events.AuthCreated{
			EventId:           utils.GenerateUUID().String(),
			Email:             email,
			VerificationToken: token.String(),
			CreatedAt:         timestamppb.New(time.Now().UTC()),
		},
	)
	if pubErr != nil {
		u.logger.Warn(
			ctx,
			"created school auth. failed to publish event",
			authIDField,
			roleField,
			eventTopicField,
			logger.NewField("error_code", ce.CodeEventPublishingFailed),
			logger.NewField("error", pubErr),
		)
		return a, at, nil
	}

	u.logger.Info(
		ctx,
		"EVENT PUBLISHED",
		authIDField,
		roleField,
		eventTopicField,
	)

	return a, at, nil
}

func (u *authUsecase) ResendVerification(ctx context.Context) (string, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ResendVerification")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return "", ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("failed to resend verification: auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)

	// Auth Fetching
	// NOTE: Resend verification is only allowed if not yet verified
	a, err := u.ar.GetByID(ctx, authCtx.AuthID)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return "", ce.NewError(
			ce.CodeAuthNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			err.Append(
				schoolIDField,
				roleField,
			).Fields()...,
		)
	}
	if err != nil {
		return "", err.Append(schoolIDField, roleField)
	}
	if a.IsVerified() {
		return "", ce.NewError(
			ce.CodeAuthAlreadyVerified,
			ce.MsgAuthAlreadyVerified,
			nil,
			authIDField,
			schoolIDField,
			roleField,
		)
	}
	if a.Email == nil {
		return "", ce.NewError(
			ce.CodeEmailNotRegistered,
			ce.MsgEmailNotRegistered,
			nil,
			authIDField,
			schoolIDField,
			roleField,
		)
	}

	// Verification Token Creation
	token := utils.GenerateUUID()
	err = u.tr.CreateVerification(
		ctx,
		&models.CreateVerificationTokenData{
			AuthID:   a.ID,
			Token:    token.String(),
			Duration: u.verificationTokenExpiry,
		},
	)
	if err != nil {
		return "", err.Append(schoolIDField, roleField)
	}

	eventTopicField := logger.NewField("event_topic", constants.EventTopicAVR)

	// Event Publishing
	pubErr := u.avrp.Publish(
		ctx,
		"auth_"+strconv.FormatInt(a.ID, 10),
		&events.AuthVerificationRequested{
			EventId:           utils.GenerateUUID().String(),
			Email:             *a.Email,
			VerificationToken: token.String(),
			CreatedAt:         timestamppb.New(time.Now().UTC()),
		},
	)
	if pubErr != nil {
		return "", ce.NewError(
			ce.CodeEventPublishingFailed,
			ce.MsgInternalServer,
			fmt.Errorf("failed to resend verification: %w", pubErr),
			authIDField,
			schoolIDField,
			roleField,
			eventTopicField,
		)
	}

	u.logger.Info(
		ctx,
		"EVENT PUBLISHED",
		authIDField,
		schoolIDField,
		roleField,
		eventTopicField,
	)

	return *a.Email, nil
}

func (u *authUsecase) VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.VerifyEmail")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("failed to verify email: auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)

	// Data Normalization
	verificationToken := strings.TrimSpace(req.VerificationToken)
	refreshToken := strings.TrimSpace(req.RefreshToken)

	// Data Validation
	if verificationToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeInvalidPayload,
			"Verification token is required",
			nil,
			authIDField,
			schoolIDField,
			roleField,
		)
	}
	if refreshToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeUnauthenticated,
			ce.MsgUnauthenticated,
			errors.New("failed to verify email: refresh token is empty"),
			authIDField,
			schoolIDField,
			roleField,
		)
	}

	// Verification Token Consumption
	authID, err := u.tr.UseVerification(ctx, verificationToken)
	if err != nil {
		return nil, nil, err.Append(authIDField, schoolIDField, roleField)
	}

	// Verification Token Ownership Validation
	// NOTE: Invalid ownership re-creates the verification token
	if authID != authCtx.AuthID {
		err := u.tr.CreateVerification(
			ctx,
			&models.CreateVerificationTokenData{
				AuthID:   authID,
				Token:    verificationToken,
				Duration: u.verificationTokenExpiry,
			},
		)
		if err != nil {
			return nil, nil, err.Append(schoolIDField, roleField)
		}
		return nil, nil, ce.NewError(
			ce.CodeTokenNotOwned,
			ce.MsgInvalidToken,
			nil,
			authIDField,
			logger.NewField("token_auth_id", authID),
			schoolIDField,
			roleField,
		)
	}

	// Verification Update
	// NOTE: Fail to update verification status re-creates the verification token
	a, err := u.ar.SetVerified(ctx, authID)
	if err != nil && err.Code() == ce.CodeAuthNotFound {
		return nil, nil, ce.NewError(
			ce.CodeAuthNotRegistered,
			ce.MsgInvalidCredentials,
			err.Unwrap(),
			err.Append(
				schoolIDField,
				roleField,
			).Fields()...,
		)
	}
	if err != nil {
		createErr := u.tr.CreateVerification(
			ctx,
			&models.CreateVerificationTokenData{
				AuthID:   authID,
				Token:    verificationToken,
				Duration: u.verificationTokenExpiry,
			},
		)
		if createErr != nil {
			return nil, nil, createErr.Append(schoolIDField, roleField)
		}
		return nil, nil, err.Append(schoolIDField, roleField)
	}

	// Session Refresh
	// NOTE: Fail to refresh session does not fail verify email process
	at, err := u.su.RefreshSession(
		ctx,
		&models.RefreshSessionReq{
			AuthID:       a.ID,
			SchoolID:     a.SchoolID,
			Role:         a.Role,
			IsVerified:   a.IsVerified(),
			RefreshToken: refreshToken,
		},
	)
	if err != nil {
		u.logger.Warn(
			ctx,
			"verified email. failed to refresh session",
			err.Append(
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Unwrap()),
			).Fields()...,
		)
		return a, nil, nil
	}

	return a, at, nil
}

func (u *authUsecase) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.IsEmailAvailable")
	defer span.End()

	// Data Normalization & Validation
	em := utils.NormalizeString(email)
	if ok, why := u.validator.Email(em); !ok {
		return false, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Email Availability Check
	return u.ar.IsEmailAvailable(ctx, em)
}
