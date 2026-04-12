package usecases

import (
	"context"
	"errors"
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
	CreateUserAuth(ctx context.Context, req *models.CreateUserAuthReq) (a *models.Auth, verificationToken *string, err *ce.Error)
	UpdateSchool(ctx context.Context, req *models.UpdateSchoolReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	ChangePassword(ctx context.Context, req *models.ChangePasswordReq) (a *models.Auth, err *ce.Error)
	ResendVerification(ctx context.Context) (email string, err *ce.Error)
	VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	RotateAuthToken(ctx context.Context, refreshToken string) (at *models.AuthToken, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
	IsUsernameAvailable(ctx context.Context, username string) (available bool, err *ce.Error)
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
			roleField,
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
	if err != nil {
		return nil, nil, err
	}

	authIDField := logger.NewField("auth_id", a.ID)
	schoolIDField := logger.NewField("school_id", a.SchoolID)
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
				schoolIDField,
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
			Role:              constants.RoleSchool,
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
			schoolIDField,
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
		schoolIDField,
		roleField,
		eventTopicField,
	)

	return a, at, nil
}

func (u *authUsecase) CreateUserAuth(ctx context.Context, req *models.CreateUserAuthReq) (*models.Auth, *string, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.CreateUserAuth")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	creatorAuthIDField := logger.NewField("creator_auth_id", authCtx.AuthID)
	creatorSchoolIDField := logger.NewField("creator_school_id", authCtx.SchoolID)
	creatorRoleField := logger.NewField("creator_role", authCtx.Role)
	creatorAuthFields := []logger.Field{
		creatorAuthIDField,
		creatorSchoolIDField,
		creatorRoleField,
	}

	// Identifier Validation
	if req.Email == nil && req.Username == nil {
		return nil, nil, ce.NewError(
			ce.CodeIdentifierNotProvided,
			ce.MsgIdentifierNotProvided,
			nil,
			creatorAuthFields...,
		)
	}

	// Data Normalization
	email := utils.NormalizeStringPtr(req.Email)
	username := utils.NormalizeStringPtr(req.Username)
	role := utils.NormalizeString(req.Role)

	// Data Validation
	if email != nil {
		if ok, why := u.validator.Email(*email); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
		}
	}
	if username != nil {
		if ok, why := u.validator.Username(*username); !ok {
			return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
		}
	}
	if ok, why := u.validator.Password(req.Password); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
	}
	if ok, why := u.validator.Role(role); !ok {
		return nil, nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
	}
	if authCtx.Role == role {
		// NOTE: Creator cannot create a new auth of the same role
		return nil, nil, ce.NewError(
			ce.CodeUnauthorizedRole,
			ce.MsgUnauthorized,
			errors.New("unable to create auth of the same role"),
			creatorAuthFields...,
		)
	}

	var a *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Email Availability Check
		if email != nil {
			available, err := u.ar.IsEmailAvailable(ctx, *email)
			if err != nil {
				return err.Append(creatorAuthFields...)
			}
			if !available {
				return ce.NewError(
					ce.CodeEmailNotAvailable,
					ce.MsgEmailNotAvailable,
					nil,
					creatorAuthFields...,
				)
			}
		}

		// Username Availability Check
		if username != nil {
			available, err := u.ar.IsUsernameAvailable(ctx, *username)
			if err != nil {
				return err.Append(creatorAuthFields...)
			}
			if !available {
				return ce.NewError(
					ce.CodeUsernameNotAvailable,
					ce.MsgUsernameNotAvailable,
					nil,
					creatorAuthFields...,
				)
			}
		}

		// Password Hashing
		hash, err := u.bcrypt.Hash(req.Password)
		if err != nil {
			return ce.NewError(
				ce.CodeBCryptHashingFailed,
				ce.MsgInternalServer,
				err,
				creatorAuthFields...,
			)
		}

		var verifiedAt *time.Time
		if email == nil {
			now := time.Now().UTC()
			verifiedAt = &now
		}

		// New Auth Creation
		na, createErr := u.ar.Create(
			ctx,
			&models.CreateAuthData{
				SchoolID:   authCtx.SchoolID,
				Email:      email,
				Username:   username,
				Password:   hash,
				Role:       role,
				VerifiedAt: verifiedAt,
			},
		)
		if createErr != nil {
			return ce.NewError(
				createErr.Code(),
				createErr.Message(),
				createErr.Unwrap(),
				creatorAuthIDField,
				creatorSchoolIDField,
				creatorRoleField,
				logger.NewField("role", role),
			)
		}

		a = na
		return nil
	})
	if err != nil && err.Code() == ce.CodeDBTransaction {
		return nil, nil, ce.NewError(
			err.Code(),
			err.Message(),
			err.Unwrap(),
			creatorAuthFields...,
		)
	}
	if err != nil {
		return nil, nil, err
	}

	// Verification Token Creation
	// NOTE: Fail to create verification token does not fail create user auth process
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
			"created user auth. failed to create verification token",
			err.Append(
				append(
					creatorAuthFields,
					logger.NewField("auth_id", a.ID),
					logger.NewField("role", a.Role),
					logger.NewField("error_code", err.Code()),
					logger.NewField("error", err.Unwrap()),
				)...,
			).Fields()...,
		)
		return a, nil, nil
	}

	verificationToken := token.String()
	return a, &verificationToken, nil
}

func (u *authUsecase) UpdateSchool(ctx context.Context, req *models.UpdateSchoolReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.UpdateSchool")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	oldSchoolIDField := logger.NewField("old_school_id", authCtx.SchoolID)
	newSchoolIDField := logger.NewField("new_school_id", req.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)
	authFields := []logger.Field{authIDField, oldSchoolIDField, newSchoolIDField, roleField}

	// Data Normalization And Validation
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeUnauthenticated,
			ce.MsgUnauthenticated,
			errors.New("refresh token is empty"),
			authFields...,
		)
	}

	var a *models.Auth
	var at *models.AuthToken
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Auth School Update
		auth, err := u.ar.UpdateSchool(ctx, authCtx.AuthID, req.SchoolID)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(
				ce.CodeAuthNotRegistered,
				ce.MsgInvalidCredentials,
				err.Unwrap(),
				err.Append(
					oldSchoolIDField,
					roleField,
				).Fields()...,
			)
		}
		if err != nil {
			return err.Append(oldSchoolIDField, roleField)
		}

		// Session Refresh
		authToken, err := u.su.RefreshSession(
			ctx,
			&models.RefreshSessionReq{
				AuthID:       auth.ID,
				SchoolID:     auth.SchoolID,
				Role:         auth.Role,
				IsVerified:   auth.IsVerified(),
				RefreshToken: refreshToken,
			},
		)
		if err != nil {
			return err
		}

		a = auth
		at = authToken
		return nil
	})
	if err != nil && err.Code() == ce.CodeDBTransaction {
		return nil, nil, ce.NewError(
			err.Code(),
			err.Message(),
			err.Unwrap(),
			authFields...,
		)
	}

	return a, at, err
}

func (u *authUsecase) ChangePassword(ctx context.Context, req *models.ChangePasswordReq) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ChangePassword")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)
	authFields := []logger.Field{authIDField, schoolIDField, roleField}

	// Data Validation
	if ok, why := u.validator.Password(req.NewPassword); !ok {
		return nil, ce.NewError(
			ce.CodeInvalidPayload,
			why,
			nil,
			authFields...,
		)
	}

	var a *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Auth Fetching
		auth, err := u.ar.GetByID(ctx, authCtx.AuthID)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(
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
			return err.Append(schoolIDField, roleField)
		}

		// Old Password Validation
		if err := u.bcrypt.Validate(auth.Password, req.OldPassword); err != nil {
			return ce.NewError(
				ce.CodeWrongPassword,
				ce.MsgInvalidCredentials,
				err,
				authFields...,
			)
		}

		// New Password Hashing
		hash, hashErr := u.bcrypt.Hash(req.NewPassword)
		if hashErr != nil {
			return ce.NewError(
				ce.CodeBCryptHashingFailed,
				ce.MsgInternalServer,
				hashErr,
				authFields...,
			)
		}

		// Password Update
		a, err = u.ar.UpdatePassword(ctx, auth.ID, hash)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(
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
			return err.Append(schoolIDField, roleField)
		}
		return nil
	})
	if err != nil && err.Code() == ce.CodeDBTransaction {
		return nil, ce.NewError(
			err.Code(),
			err.Message(),
			err.Unwrap(),
			authFields...,
		)
	}

	return a, err
}

func (u *authUsecase) ResendVerification(ctx context.Context) (string, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.ResendVerification")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return "", ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)
	authFields := []logger.Field{authIDField, schoolIDField, roleField}

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
			authFields...,
		)
	}
	if a.Email == nil {
		return "", ce.NewError(
			ce.CodeEmailNotRegistered,
			ce.MsgEmailNotRegistered,
			nil,
			authFields...,
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
			pubErr,
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
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)
	authFields := []logger.Field{authIDField, schoolIDField, roleField}

	// Data Normalization
	verificationToken := strings.TrimSpace(req.VerificationToken)
	refreshToken := strings.TrimSpace(req.RefreshToken)

	// Data Validation
	if verificationToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeInvalidPayload,
			"Verification token is required",
			nil,
			authFields...,
		)
	}
	if refreshToken == "" {
		return nil, nil, ce.NewError(
			ce.CodeUnauthenticated,
			ce.MsgUnauthenticated,
			errors.New("refresh token is empty"),
			authFields...,
		)
	}

	// Verification Token Consumption
	authID, err := u.tr.UseVerification(ctx, verificationToken)
	if err != nil {
		return nil, nil, err.Append(authFields...)
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

func (u *authUsecase) RotateAuthToken(ctx context.Context, refreshToken string) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.RotateAuthToken")
	defer span.End()

	// Data Normalization And Validation
	token := strings.TrimSpace(refreshToken)
	if token == "" {
		return nil, ce.NewError(
			ce.CodeUnauthenticated,
			ce.MsgUnauthenticated,
			errors.New("refresh token is empty"),
		)
	}

	var at *models.AuthToken
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Session Fetching
		s, err := u.su.GetSession(ctx, token)
		if err != nil {
			return err
		}

		authIDField := logger.NewField("auth_id", s.AuthID)

		// Session Expiration Check
		if s.ExpiresAt.Before(time.Now().UTC()) {
			return ce.NewError(
				ce.CodeSessionExpired,
				ce.MsgSessionExpired,
				nil,
				authIDField,
			)
		}

		// Auth Fetching
		a, err := u.ar.GetByID(ctx, s.AuthID)
		if err != nil && err.Code() == ce.CodeAuthNotFound {
			return ce.NewError(
				ce.CodeAuthNotRegistered,
				ce.MsgInvalidCredentials,
				err.Unwrap(),
				err.Fields()...,
			)
		}
		if err != nil {
			return err
		}

		// Session Refresh
		at, err = u.su.RefreshSession(
			ctx,
			&models.RefreshSessionReq{
				AuthID:       a.ID,
				SchoolID:     a.SchoolID,
				Role:         a.Role,
				IsVerified:   a.IsVerified(),
				RefreshToken: token,
			},
		)
		return err
	})

	return at, err
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

func (u *authUsecase) IsUsernameAvailable(ctx context.Context, username string) (bool, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.IsUsernameAvailable")
	defer span.End()

	// Data Normalization & Validation
	un := utils.NormalizeString(username)
	if ok, why := u.validator.Username(un); !ok {
		return false, ce.NewError(ce.CodeInvalidPayload, why, nil)
	}

	// Username Availability Check
	return u.ar.IsUsernameAvailable(ctx, un)
}
