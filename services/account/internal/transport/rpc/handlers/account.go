package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/account/internal/models"
	"github.com/ritchieridanko/klasshub/services/account/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AccountHandler struct {
	apis.UnimplementedAccountServiceServer
	au usecases.AccountUsecase
}

func NewAccountHandler(au usecases.AccountUsecase) *AccountHandler {
	return &AccountHandler{au: au}
}

func (h *AccountHandler) CreateSchoolProfile(ctx context.Context, req *apis.CreateSchoolProfileRequest) (*apis.CreateSchoolProfileResponse, error) {
	s, a, at, err := h.au.CreateSchoolProfile(
		ctx,
		&models.CreateSchoolProfileReq{
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
			RefreshToken:  req.GetRefreshToken(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateSchoolProfileResponse{
		School:    h.toSchool(s),
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AccountHandler) toSchool(s *models.School) *apis.School {
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
		CreatedAt:      utils.ToTimestamp(s.CreatedAt),
	}
}

func (h *AccountHandler) toAuth(a *models.Auth) *apis.Auth {
	if a == nil {
		return nil
	}
	return &apis.Auth{
		Email:             a.Email,
		Username:          a.Username,
		Role:              a.Role,
		IsVerified:        a.IsVerified,
		SchoolExists:      a.SchoolExists,
		PasswordChangedAt: utils.ToTimestamp(a.PasswordChangedAt),
	}
}

func (h *AccountHandler) toAuthToken(at *models.AuthToken) *apis.AuthToken {
	if at == nil {
		return nil
	}
	return &apis.AuthToken{
		AccessToken:  h.toAccessToken(at.AccessToken),
		RefreshToken: h.toRefreshToken(at.RefreshToken),
	}
}

func (h *AccountHandler) toAccessToken(at *models.AccessToken) *apis.AccessToken {
	if at == nil {
		return nil
	}
	return &apis.AccessToken{
		Token:     at.Token,
		ExpiresIn: at.ExpiresIn,
	}
}

func (h *AccountHandler) toRefreshToken(rt *models.RefreshToken) *apis.RefreshToken {
	if rt == nil {
		return nil
	}
	return &apis.RefreshToken{
		Token:     rt.Token,
		ExpiresIn: rt.ExpiresIn,
	}
}
