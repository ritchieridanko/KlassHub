package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/models"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type SchoolClient interface {
	CreateSchool(ctx context.Context, req *models.CreateSchoolReq) (schoolID int64, s *models.School, err *ce.Error)
}

type schoolClient struct {
	client apis.SchoolServiceClient
}

func NewSchoolClient(c apis.SchoolServiceClient) SchoolClient {
	return &schoolClient{client: c}
}

func (c *schoolClient) CreateSchool(ctx context.Context, req *models.CreateSchoolReq) (int64, *models.School, *ce.Error) {
	resp, err := c.client.CreateSchool(
		ctx,
		&apis.CreateSchoolRequest{
			NPSN:          req.NPSN,
			Name:          req.Name,
			Level:         req.Level,
			Ownership:     req.Ownership,
			Accreditation: req.Accreditation,
			EstablishedAt: utils.ToTimestamp(req.EstablishedAt),
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
		return 0, nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "school"),
		)
	}
	return resp.GetSchoolId(), c.toSchool(resp.GetSchool()), nil
}

func (c *schoolClient) toSchool(s *apis.School) *models.School {
	if s == nil {
		return nil
	}
	return &models.School{
		NPSN:           s.NPSN,
		NPSNVerifiedAt: utils.ToTime(s.GetNPSNVerifiedAt()),
		Name:           s.GetName(),
		Level:          s.GetLevel(),
		Ownership:      s.GetOwnership(),
		ProfilePicture: s.ProfilePicture,
		ProfileBanner:  s.ProfileBanner,
		Accreditation:  s.Accreditation,
		EstablishedAt:  utils.ToTime(s.GetEstablishedAt()),
		Province:       s.GetProvince(),
		CityRegency:    s.GetCityRegency(),
		District:       s.GetDistrict(),
		Subdistrict:    s.GetSubdistrict(),
		Street:         s.GetStreet(),
		Postcode:       s.GetPostcode(),
		Phone:          s.Phone,
		Email:          s.Email,
		Website:        s.Website,
		Timezone:       s.GetTimezone(),
		CreatedAt:      utils.ToTime(s.GetCreatedAt()),
	}
}
