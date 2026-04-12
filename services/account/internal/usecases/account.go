package usecases

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ritchieridanko/klasshub/services/account/internal/clients"
	"github.com/ritchieridanko/klasshub/services/account/internal/constants"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/account/internal/models"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/metadata"
	"github.com/ritchieridanko/klasshub/shared/contract/events/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountUsecase interface {
	CreateSchoolProfile(ctx context.Context, req *models.CreateSchoolProfileReq) (s *models.School, a *models.Auth, at *models.AuthToken, err *ce.Error)
	CreateUserAccount(ctx context.Context, req *models.CreateUserAccountReq) (a *models.Auth, u *models.User, err *ce.Error)
}

type accountUsecase struct {
	appName string
	ac      clients.AuthClient
	sc      clients.SchoolClient
	uc      clients.UserClient
	acp     *publisher.Publisher
	asufp   *publisher.Publisher
	ucfp    *publisher.Publisher
	logger  *logger.Logger
}

func NewAccountUsecase(appName string, ac clients.AuthClient, sc clients.SchoolClient, uc clients.UserClient, acp *publisher.Publisher, asufp *publisher.Publisher, ucfp *publisher.Publisher, l *logger.Logger) AccountUsecase {
	return &accountUsecase{
		appName: appName,
		ac:      ac,
		sc:      sc,
		uc:      uc,
		acp:     acp,
		asufp:   asufp,
		ucfp:    ucfp,
		logger:  l,
	}
}

func (u *accountUsecase) CreateSchoolProfile(ctx context.Context, req *models.CreateSchoolProfileReq) (*models.School, *models.Auth, *models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "account.usecase.CreateSchoolProfile")
	defer span.End()

	authCtx := utils.CtxAuth(ctx)
	if authCtx == nil {
		return nil, nil, nil, ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		)
	}

	authIDField := logger.NewField("auth_id", authCtx.AuthID)
	schoolIDField := logger.NewField("school_id", authCtx.SchoolID)
	roleField := logger.NewField("role", authCtx.Role)
	authFields := []logger.Field{authIDField, schoolIDField, roleField}

	// School Creation
	schoolID, s, err := u.sc.CreateSchool(
		metadata.ToOutgoingCtx(
			ctx,
			metadata.Auth(authCtx, true, false, true, true)...,
		),
		&models.CreateSchoolReq{
			NPSN:          req.NPSN,
			Name:          req.Name,
			Level:         req.Level,
			Ownership:     req.Ownership,
			Accreditation: req.Accreditation,
			EstablishedAt: req.EstablishedAt,
			Province:      req.Province,
			CityRegency:   req.CityRegency,
			District:      req.District,
			Subdistrict:   req.Subdistrict,
			Street:        req.Street,
			Postcode:      req.Postcode,
			Phone:         req.Phone,
			Email:         req.Email,
			Website:       req.Website,
			Timezone:      req.Timezone,
		},
	)
	if err != nil {
		return nil, nil, nil, err.Append(authFields...)
	}

	// Auth School Update
	a, at, err := u.ac.UpdateSchool(
		metadata.ToOutgoingCtx(
			ctx,
			metadata.Auth(authCtx, true, true, true, true)...,
		),
		&models.UpdateSchoolReq{
			SchoolID:     schoolID,
			RefreshToken: req.RefreshToken,
		},
	)
	if err != nil {
		newSchoolIDField := logger.NewField("new_school_id", schoolID)
		eventTopicField := logger.NewField("event_topic", constants.EventTopicASUF)

		// Created School Cancellation (Async)
		pubErr := u.asufp.Publish(
			ctx,
			"school_"+strconv.FormatInt(schoolID, 10),
			&events.AuthSchoolUpdateFailed{
				EventId:   utils.GenerateUUID().String(),
				SchoolId:  schoolID,
				CreatedAt: timestamppb.New(time.Now().UTC()),
			},
		)
		if pubErr != nil {
			u.logger.Warn(
				ctx,
				"failed to create school profile. failed to publish event",
				authIDField,
				schoolIDField,
				newSchoolIDField,
				roleField,
				eventTopicField,
				logger.NewField("error_code", ce.CodeEventPublishingFailed),
				logger.NewField("error", pubErr),
			)
		} else {
			u.logger.Info(
				ctx,
				"EVENT PUBLISHED",
				authIDField,
				schoolIDField,
				newSchoolIDField,
				roleField,
				eventTopicField,
			)
		}

		return nil, nil, nil, err.Append(
			authIDField,
			schoolIDField,
			newSchoolIDField,
			roleField,
		)
	}

	return s, a, at, nil
}

func (u *accountUsecase) CreateUserAccount(ctx context.Context, req *models.CreateUserAccountReq) (*models.Auth, *models.User, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "account.usecase.CreateUserAccount")
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

	// School Existence Check
	exists, err := u.sc.SchoolExists(ctx, authCtx.SchoolID)
	if err != nil {
		return nil, nil, err.Append(creatorAuthFields...)
	}
	if !exists {
		return nil, nil, ce.NewError(
			ce.CodeSchoolNotRegistered,
			ce.MsgInvalidCredentials,
			nil,
			creatorAuthFields...,
		)
	}

	// Auth Creation
	authID, schoolID, a, verificationToken, err := u.ac.CreateUserAuth(
		metadata.ToOutgoingCtx(
			ctx,
			metadata.Auth(authCtx, true, true, true, true)...,
		),
		&models.CreateUserAuthReq{
			Email:    req.Email,
			Username: req.Username,
			Password: req.Password,
			Role:     req.Role,
		},
	)
	if err != nil {
		return nil, nil, err.Append(creatorAuthFields...)
	}

	authFields := append(
		creatorAuthFields,
		logger.NewField("auth_id", authID),
		logger.NewField("role", a.Role),
	)

	// User Creation
	user, err := u.uc.CreateUser(
		metadata.ToOutgoingCtx(
			ctx,
			metadata.Auth(authCtx, true, true, true, true)...,
		),
		&models.CreateUserReq{
			AuthID:       authID,
			SchoolID:     schoolID,
			SchoolUserID: req.SchoolUserID,
			Role:         a.Role,
			Name:         req.Name,
			Birthplace:   req.Birthplace,
			Birthdate:    req.Birthdate,
			Sex:          req.Sex,
		},
	)
	if err != nil {
		fields := append(
			authFields,
			logger.NewField("event_topic", constants.EventTopicUCF),
		)

		// Created Auth Cancellation (Async)
		pubErr := u.ucfp.Publish(
			ctx,
			"auth_"+strconv.FormatInt(authID, 10),
			&events.UserCreationFailed{
				EventId:   utils.GenerateUUID().String(),
				AuthId:    authID,
				CreatedAt: timestamppb.New(time.Now().UTC()),
			},
		)
		if pubErr != nil {
			u.logger.Warn(
				ctx,
				"failed to create user account. failed to publish event",
				append(
					fields,
					logger.NewField("error_code", ce.CodeEventPublishingFailed),
					logger.NewField("error", pubErr),
				)...,
			)
		} else {
			u.logger.Info(ctx, "EVENT PUBLISHED", fields...)
		}

		return nil, nil, err.Append(fields...)
	}

	// Verification Requirement Check
	if a.Email != nil && verificationToken != nil {
		fields := append(
			authFields,
			logger.NewField("event_topic", constants.EventTopicAC),
		)

		// Event Publishing
		// NOTE: Fail to publish event does not fail create user account process
		err := u.acp.Publish(
			ctx,
			"auth_"+strconv.FormatInt(authID, 10),
			&events.AuthCreated{
				EventId:           utils.GenerateUUID().String(),
				Role:              a.Role,
				Email:             *a.Email,
				VerificationToken: *verificationToken,
				CreatedAt:         timestamppb.New(time.Now().UTC()),
			},
		)
		if err != nil {
			u.logger.Warn(
				ctx,
				"created user account. failed to publish event",
				append(
					fields,
					logger.NewField("error_code", ce.CodeEventPublishingFailed),
					logger.NewField("error", err),
				)...,
			)
		} else {
			u.logger.Info(ctx, "EVENT PUBLISHED", fields...)
		}
	}

	return a, user, nil
}
