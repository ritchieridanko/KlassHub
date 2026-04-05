package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/ritchieridanko/klasshub/services/school/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/validator"
	"go.opentelemetry.io/otel"
)

type SchoolUsecase interface {
	CreateSchool(ctx context.Context, req *models.CreateSchoolReq) (s *models.School, err *ce.Error)
}

type schoolUsecase struct {
	appName   string
	sr        repositories.SchoolRepository
	validator *validator.Validator
}

func NewSchoolUsecase(appName string, sr repositories.SchoolRepository, v *validator.Validator) SchoolUsecase {
	return &schoolUsecase{
		appName:   appName,
		sr:        sr,
		validator: v,
	}
}

func (u *schoolUsecase) CreateSchool(ctx context.Context, req *models.CreateSchoolReq) (*models.School, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "school.usecase.CreateSchool")
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
	roleField := logger.NewField("role", authCtx.Role)
	authFields := []logger.Field{authIDField, roleField}

	// Data Normalization
	npsn := utils.NormalizeStringPtr(req.NPSN)
	name := strings.TrimSpace(req.Name)
	level := utils.NormalizeString(req.Level)
	ownership := utils.NormalizeString(req.Ownership)
	accreditation := utils.NormalizeStringPtr(req.Accreditation)
	province := utils.NormalizeString(req.Province)
	cityRegency := utils.NormalizeString(req.CityRegency)
	district := utils.NormalizeString(req.District)
	subdistrict := utils.NormalizeString(req.Subdistrict)
	street := strings.TrimSpace(req.Street)
	postcode := utils.NormalizeString(req.Postcode)
	phone := utils.NormalizeStringPtr(req.Phone)
	email := utils.NormalizeStringPtr(req.Email)
	website := utils.NormalizeStringPtr(req.Website)
	timezone := utils.NormalizeString(req.Timezone)

	// Data Validation
	if npsn != nil {
		if ok, why := u.validator.NPSN(*npsn); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if ok, why := u.validator.SchoolName(name); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.SchoolLevel(level); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.SchoolOwnership(ownership); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if accreditation != nil {
		if ok, why := u.validator.SchoolAccreditation(*accreditation); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if req.EstablishedAt != nil {
		if ok, why := u.validator.SchoolEstablishedAt(*req.EstablishedAt); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if ok, why := u.validator.Street(street); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if ok, why := u.validator.Postcode(postcode); !ok {
		return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
	}
	if phone != nil {
		if ok, why := u.validator.Phone(*phone); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if email != nil {
		if ok, why := u.validator.Email(*email); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}
	if website != nil {
		if ok, why := u.validator.URL(*website); !ok {
			return nil, ce.NewError(ce.CodeInvalidPayload, why, nil, authFields...)
		}
	}

	// School Creation
	s, err := u.sr.Create(
		ctx,
		&models.CreateSchoolData{
			NPSN:          npsn,
			Name:          name,
			Level:         level,
			Ownership:     ownership,
			Accreditation: accreditation,
			EstablishedAt: req.EstablishedAt,
			Province:      province,
			CityRegency:   cityRegency,
			District:      district,
			Subdistrict:   subdistrict,
			Street:        street,
			Postcode:      postcode,
			Phone:         phone,
			Email:         email,
			Website:       website,
			Timezone:      timezone,
		},
	)
	if err != nil {
		return nil, err.Append(authFields...)
	}

	return s, nil
}
