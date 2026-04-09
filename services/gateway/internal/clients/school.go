package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SchoolClient interface {
	GetMe(ctx context.Context) (s *models.School, err *ce.Error)
}

type schoolClient struct {
	client apis.SchoolServiceClient
}

func NewSchoolClient(c apis.SchoolServiceClient) SchoolClient {
	return &schoolClient{client: c}
}

func (c *schoolClient) GetMe(ctx context.Context) (*models.School, *ce.Error) {
	resp, err := c.client.GetMe(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "school"),
		)
	}
	return c.toSchool(resp.GetSchool()), nil
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
