package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/user/internal/constants"
	"github.com/ritchieridanko/klasshub/services/user/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/validator"
	"go.opentelemetry.io/otel"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *models.CreateUserReq) (u *models.User, err *ce.Error)
	GetMe(ctx context.Context) (u *models.User, err *ce.Error)
}

type userUsecase struct {
	appName   string
	ur        repositories.UserRepository
	validator *validator.Validator
}

func NewUserUsecase(appName string, ur repositories.UserRepository, v *validator.Validator) UserUsecase {
	return &userUsecase{
		appName:   appName,
		ur:        ur,
		validator: v,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *models.CreateUserReq) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "user.usecase.CreateUser")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", req.AuthID)
	creatorAuthIDField := logger.NewField("creator_auth_id", authCtx.AuthID)
	creatorSchoolIDField := logger.NewField("creator_school_id", authCtx.SchoolID)
	creatorRoleField := logger.NewField("creator_role", authCtx.Role)
	authFields := []logger.Field{
		authIDField,
		creatorAuthIDField,
		creatorSchoolIDField,
		creatorRoleField,
	}

	var createdBy *uuid.UUID
	if authCtx.Role == constants.RoleAdministrator {
		// Creator And New User Role Validation
		// NOTE: Administrators cannot create new users of type administrator
		if req.Role == constants.RoleAdministrator {
			return nil, ce.NewError(
				ce.CodeUnauthorizedRole,
				ce.MsgUnauthorized,
				errors.New("unable to create administrator user"),
				authFields...,
			)
		}

		// Creator User Fetching
		user, err := u.ur.GetByAuthID(ctx, authCtx.AuthID)
		if err != nil && err.Code() == ce.CodeUserNotFound {
			return nil, ce.NewError(
				ce.CodeUserNotRegistered,
				ce.MsgInvalidCredentials,
				err.Unwrap(),
				authFields...,
			)
		}
		if err != nil {
			return nil, ce.NewError(
				err.Code(),
				err.Message(),
				err.Unwrap(),
				authFields...,
			)
		}

		createdBy = &user.ID
	}

	// Data Normalization
	schoolUserID := utils.TrimSpacePtr(req.SchoolUserID)
	name := strings.TrimSpace(req.Name)
	birthplace := utils.NormalizeString(req.Birthplace)
	sex := utils.NormalizeString(req.Sex)

	// Data Validation
	if schoolUserID != nil {
		if ok, why := u.validator.SchoolUserID(*schoolUserID); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if ok, why := u.validator.Role(req.Role); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Name(name); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Birthplace(birthplace); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Birthdate(req.Birthdate); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Sex(sex); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}

	// New User ID Generation
	uuid, err := utils.GenerateUUIDv7()
	if err != nil {
		return nil, ce.NewError(
			ce.CodeUUIDGenerationFailed,
			ce.MsgInternalServer,
			err,
			authFields...,
		)
	}

	// New User Creation
	nu, createErr := u.ur.Create(
		ctx,
		&models.CreateUserData{
			ID:           uuid,
			AuthID:       req.AuthID,
			SchoolID:     req.SchoolID,
			SchoolUserID: schoolUserID,
			Role:         req.Role,
			Name:         name,
			Birthplace:   birthplace,
			Birthdate:    req.Birthdate,
			Sex:          sex,
			CreatedBy:    createdBy,
		},
	)
	if createErr != nil {
		return nil, createErr.Append(
			creatorAuthIDField,
			creatorSchoolIDField,
			creatorRoleField,
		)
	}

	return nu, nil
}

func (u *userUsecase) GetMe(ctx context.Context) (*models.User, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "user.usecase.GetMe")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	// User Fetching
	user, err := u.ur.GetByAuthID(ctx, authCtx.AuthID)
	if err != nil {
		return nil, err.Append(
			logger.NewField("school_id", authCtx.SchoolID),
			logger.NewField("role", authCtx.Role),
		)
	}
	return user, nil
}
