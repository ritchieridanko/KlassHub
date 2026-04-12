package clients

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AccountClient interface {
	CreateSchoolProfile(ctx context.Context, req *models.CreateSchoolProfileReq) (s *models.School, a *models.Auth, at *models.AuthToken, err *ce.Error)
	CreateUserAccount(ctx context.Context, req *models.CreateUserAccountReq) (a *models.Auth, u *models.User, err *ce.Error)
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

func (c *accountClient) CreateUserAccount(ctx context.Context, req *models.CreateUserAccountReq) (*models.Auth, *models.User, *ce.Error) {
	resp, err := c.client.CreateUserAccount(
		ctx,
		&apis.CreateUserAccountRequest{
			// Auth
			Email:    req.Email,
			Username: req.Username,
			Password: req.Password,
			Role:     req.Role,

			// User
			SchoolUserId: req.SchoolUserID,
			Name:         req.Name,
			Birthplace:   req.Birthplace,
			Birthdate:    utils.ToTimestamp(req.Birthdate),
			Sex:          req.Sex,
		},
	)
	if err != nil {
		return nil, nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "account"),
		)
	}
	return c.toAuthFromAA(resp.GetAuth()),
		c.toUserFromUA(resp.GetUser()),
		nil
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

func (c *accountClient) toAuthFromAA(aa *apis.AuthAdmin) *models.Auth {
	if aa == nil {
		return nil
	}
	return &models.Auth{
		Email:      aa.Email,
		Username:   aa.Username,
		Role:       aa.GetRole(),
		IsVerified: aa.GetIsVerified(),
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

func (c *accountClient) toUserFromUA(ua *apis.UserAdmin) *models.User {
	if ua == nil {
		return nil
	}

	var createdBy *uuid.UUID
	if ua.CreatedBy != nil {
		creator := utils.ToUUIDMust(ua.GetCreatedBy())
		createdBy = &creator
	}

	return &models.User{
		ID:             utils.ToUUIDMust(ua.GetId()),
		SchoolUserID:   ua.SchoolUserId,
		Role:           ua.GetRole(),
		Name:           ua.GetName(),
		Nickname:       ua.Nickname,
		Birthplace:     ua.GetBirthplace(),
		Birthdate:      utils.ToTime(ua.GetBirthdate()),
		Sex:            ua.GetSex(),
		Phone:          ua.Phone,
		ProfilePicture: ua.ProfilePicture,
		ProfileBanner:  ua.ProfileBanner,
		CreatedBy:      createdBy,
		CreatedAt:      utils.ToTime(ua.GetCreatedAt()),
		UpdatedAt:      utils.ToTime(ua.GetUpdatedAt()),
	}
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
