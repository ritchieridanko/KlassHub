package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type SchoolHandler struct {
	apis.UnimplementedSchoolServiceServer
	su usecases.SchoolUsecase
}

func NewSchoolHandler(su usecases.SchoolUsecase) *SchoolHandler {
	return &SchoolHandler{su: su}
}

func (h *SchoolHandler) CreateSchool(ctx context.Context, req *apis.CreateSchoolRequest) (*apis.CreateSchoolResponse, error) {
	s, err := h.su.CreateSchool(
		ctx,
		&models.CreateSchoolReq{
			NPSN:          req.NPSN,
			Name:          req.GetName(),
			Level:         req.GetLevel(),
			Ownership:     req.GetOwnership(),
			Accreditation: req.Accreditation,
			EstablishedAt: utils.ToTime(req.GetEstablishedAt()),
			Province:      req.GetProvince(),
			CityRegency:   req.GetCityRegency(),
			District:      req.GetDistrict(),
			Subdistrict:   req.GetSubdistrict(),
			Street:        req.GetStreet(),
			Postcode:      req.GetPostcode(),
			Phone:         req.Phone,
			Email:         req.Email,
			Website:       req.Website,
			Timezone:      req.GetTimezone(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateSchoolResponse{
		SchoolId: s.ID,
		School:   h.toSchool(s),
	}, nil
}

func (h *SchoolHandler) toSchool(s *models.School) *apis.School {
	if s == nil {
		return nil
	}
	return &apis.School{
		NPSN:           s.NPSN,
		NPSNVerifiedAt: utils.ToTimestamp(s.NPSNVerifiedAt),
		Name:           s.Name,
		Level:          s.Level,
		Ownership:      s.Ownership,
		ProfilePicture: s.ProfilePicture,
		ProfileBanner:  s.ProfileBanner,
		Accreditation:  s.Accreditation,
		EstablishedAt:  utils.ToTimestamp(s.EstablishedAt),
		Province:       s.Province,
		CityRegency:    s.CityRegency,
		District:       s.District,
		Subdistrict:    s.Subdistrict,
		Street:         s.Street,
		Postcode:       s.Postcode,
		Phone:          s.Phone,
		Email:          s.Email,
		Website:        s.Website,
		Timezone:       s.Timezone,
		CreatedAt:      utils.ToTimestamp(&s.CreatedAt),
	}
}
