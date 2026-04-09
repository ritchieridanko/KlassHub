package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AccountClient interface {
	CreateSchoolProfile(ctx context.Context, req *models.CreateSchoolProfileReq) (s *models.School, a *models.Auth, at *models.AuthToken, err *ce.Error)
}

type accountClient struct {
	client apis.AccountServiceClient
}

func NewAccountClient(c apis.AccountServiceClient) AccountClient {
	return &accountClient{client: c}
}

func (c *accountClient) CreateSchoolProfile(ctx context.Context, req *models.CreateSchoolProfileReq) (*models.School, *models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.CreateSchoolProfile(
		ctx,
		&apis.CreateSchoolProfileRequest{
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
			RefreshToken:  req.RefreshToken,
		},
	)
	if err != nil {
		return nil, nil, nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "account"),
		)
	}
	return c.toSchool(resp.GetSchool()),
		c.toAuth(resp.GetAuth()),
		c.toAuthToken(resp.GetAuthToken()),
		nil
}

func (c *accountClient) toSchool(s *apis.School) *models.School {
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

func (c *accountClient) toAuth(a *apis.Auth) *models.Auth {
	if a == nil {
		return nil
	}
	return &models.Auth{
		Email:             a.Email,
		Username:          a.Username,
		Role:              a.GetRole(),
		IsVerified:        a.GetIsVerified(),
		SchoolExists:      a.GetSchoolExists(),
		PasswordChangedAt: utils.ToTime(a.GetPasswordChangedAt()),
	}
}

func (c *accountClient) toAuthToken(at *apis.AuthToken) *models.AuthToken {
	if at == nil {
		return nil
	}
	return &models.AuthToken{
		AccessToken:  c.toAccessToken(at.GetAccessToken()),
		RefreshToken: c.toRefreshToken(at.GetRefreshToken()),
	}
}

func (c *accountClient) toAccessToken(at *apis.AccessToken) *models.AccessToken {
	if at == nil {
		return nil
	}
	return &models.AccessToken{
		Token:     at.GetToken(),
		ExpiresIn: at.GetExpiresIn(),
	}
}

func (c *accountClient) toRefreshToken(rt *apis.RefreshToken) *models.RefreshToken {
	if rt == nil {
		return nil
	}
	return &models.RefreshToken{
		Token:     rt.GetToken(),
		ExpiresIn: rt.GetExpiresIn(),
	}
}
