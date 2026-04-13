package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/ritchieridanko/klasshub/services/course/internal/clients"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/course/internal/models"
	"github.com/ritchieridanko/klasshub/services/course/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils/validator"
	"go.opentelemetry.io/otel"
)

type CourseUsecase interface {
	CreateCourse(ctx context.Context, req *models.CreateCourseReq) (c *models.Course, err *ce.Error)
}

type courseUsecase struct {
	appName   string
	sc        clients.SchoolClient
	cr        repositories.CourseRepository
	validator *validator.Validator
}

func NewCourseUsecase(appName string, sc clients.SchoolClient, cr repositories.CourseRepository, v *validator.Validator) CourseUsecase {
	return &courseUsecase{
		appName:   appName,
		sc:        sc,
		cr:        cr,
		validator: v,
	}
}

func (u *courseUsecase) CreateCourse(ctx context.Context, req *models.CreateCourseReq) (*models.Course, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "course.usecase.CreateCourse")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, ce.NewError(
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

	// School Existence Check
	exists, err := u.sc.SchoolExists(ctx, authCtx.SchoolID)
	if err != nil {
		return nil, err.Append(creatorAuthFields...)
	}
	if !exists {
		return nil, ce.NewError(
			ce.CodeSchoolNotRegistered,
			ce.MsgInvalidCredentials,
			nil,
			creatorAuthFields...,
		)
	}

	// Data Normalization
	schoolCourseID := utils.TrimSpacePtr(req.SchoolCourseID)
	name := strings.TrimSpace(req.Name)

	// Data Validation
	if schoolCourseID != nil {
		if ok, why := u.validator.SchoolCourseID(*schoolCourseID); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
		}
	}
	if ok, why := u.validator.CourseName(name); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
	}
	if req.Description != nil {
		if ok, why := u.validator.CourseDesc(*req.Description); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
		}
	}
	if req.CoursePicture != nil {
		if ok, why := u.validator.URL(*req.CoursePicture); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, creatorAuthFields...)
		}
	}

	// Course ID Generation
	uuid, genErr := utils.GenerateUUIDv7()
	if genErr != nil {
		return nil, ce.NewError(
			ce.CodeUUIDGenerationFailed,
			ce.MsgInternalServer,
			genErr,
			creatorAuthFields...,
		)
	}

	// Course Creation
	c, err := u.cr.Create(
		ctx,
		&models.CreateCourseData{
			ID:             uuid,
			SchoolID:       authCtx.SchoolID,
			SchoolCourseID: schoolCourseID,
			Name:           name,
			Description:    req.Description,
			CoursePicture:  req.CoursePicture,
		},
	)
	if err != nil {
		return nil, err.Append(creatorAuthIDField, creatorRoleField)
	}

	return c, nil
}
