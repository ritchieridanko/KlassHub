package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/clients"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/cookie"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/metadata"
)

type AccountHandler struct {
	acc    clients.AccountClient
	cookie *cookie.Cookie
}

func NewAccountHandler(acc clients.AccountClient, c *cookie.Cookie) *AccountHandler {
	return &AccountHandler{acc: acc, cookie: c}
}

func (h *AccountHandler) CreateSchoolProfile(ctx *gin.Context) {
	var payload dtos.CreateSchoolProfileRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	refreshToken, err := ctx.Cookie(constants.CookieKeyRefreshToken)
	if errors.Is(err, ce.ErrCookieNotFound) {
		ce.NewError(ce.CodeRefreshTokenNotFound, ce.MsgInvalidSession, err).Bind(ctx)
		return
	}
	if err != nil {
		ce.NewError(ce.CodeInternal, ce.MsgInternalServer, err).Bind(ctx)
		return
	}

	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		ce.NewError(ce.CodeRefreshTokenNotFound, ce.MsgInvalidSession, nil).Bind(ctx)
		return
	}

	s, a, at, createErr := h.acc.CreateSchoolProfile(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, true)...,
		),
		&models.CreateSchoolProfileReq{
			NPSN:          payload.NPSN,
			Name:          payload.Name,
			Level:         payload.Level,
			Ownership:     payload.Ownership,
			Accreditation: payload.Accreditation,
			EstablishedAt: payload.EstablishedAt,
			Province:      payload.Province,
			CityRegency:   payload.CityRegency,
			District:      payload.District,
			Subdistrict:   payload.Subdistrict,
			Street:        payload.Street,
			Postcode:      payload.Postcode,
			Phone:         payload.Phone,
			Email:         payload.Email,
			Website:       payload.Website,
			Timezone:      payload.Timezone,
			RefreshToken:  refreshToken,
		},
	)
	if createErr != nil {
		createErr.Bind(ctx)
		return
	}

	if at != nil && at.RefreshToken != nil {
		h.cookie.Set(
			ctx,
			constants.CookieKeyRefreshToken,
			at.RefreshToken.Token,
			"/",
			int(at.RefreshToken.ExpiresIn),
		)
	}

	utils.SetResponse(
		ctx,
		http.StatusCreated,
		"School created successfully",
		dtos.CreateSchoolProfileResponse{
			School:      h.toSchool(s),
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
		},
	)
}

func (h *AccountHandler) CreateUserAccount(ctx *gin.Context) {
	var payload dtos.CreateUserAccountRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	authCtx := utils.CtxAuth(ctx.Request.Context())
	if authCtx == nil {
		ce.NewError(
			ce.CodeMissingContextValue,
			ce.MsgInternalServer,
			errors.New("auth missing from context"),
		).Bind(
			ctx,
		)
		return
	}

	a, u, err := h.acc.CreateUserAccount(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, true)...,
		),
		&models.CreateUserAccountReq{
			// Auth
			Email:    payload.Email,
			Username: payload.Username,
			Password: payload.Password,
			Role:     payload.Role,

			// User
			SchoolUserID: payload.SchoolUserID,
			Name:         payload.Name,
			Birthplace:   payload.Birthplace,
			Birthdate:    payload.Birthdate,
			Sex:          payload.Sex,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusCreated,
		"User created successfully",
		dtos.CreateUserAccountResponse{
			Auth: h.toAuthAdmin(a),
			User: h.toUserAdmin(u),
		},
	)
}

func (h *AccountHandler) toAuth(a *models.Auth) *dtos.Auth {
	if a == nil {
		return nil
	}
	return &dtos.Auth{
		Email:             a.Email,
		Username:          a.Username,
		Role:              a.Role,
		IsVerified:        a.IsVerified,
		SchoolExists:      a.SchoolExists,
		PasswordChangedAt: a.PasswordChangedAt,
	}
}

func (h *AccountHandler) toAuthAdmin(a *models.Auth) *dtos.AuthAdmin {
	if a == nil {
		return nil
	}
	return &dtos.AuthAdmin{
		Email:      a.Email,
		Username:   a.Username,
		Role:       a.Role,
		IsVerified: a.IsVerified,
	}
}

func (h *AccountHandler) toAccessToken(at *models.AuthToken) *dtos.AccessToken {
	if at == nil || at.AccessToken == nil {
		return nil
	}
	return &dtos.AccessToken{
		Token:     at.AccessToken.Token,
		ExpiresIn: at.AccessToken.ExpiresIn,
	}
}

func (h *AccountHandler) toUserAdmin(u *models.User) *dtos.UserAdmin {
	if u == nil {
		return nil
	}

	var createdBy *string
	if u.CreatedBy != nil {
		creator := u.CreatedBy.String()
		createdBy = &creator
	}

	return &dtos.UserAdmin{
		ID:             u.ID.String(),
		SchoolUserID:   u.SchoolUserID,
		Role:           u.Role,
		Name:           u.Name,
		Birthplace:     u.Birthplace,
		Birthdate:      u.Birthdate,
		Sex:            u.Sex,
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		CreatedBy:      createdBy,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func (h *AccountHandler) toSchool(s *models.School) *dtos.School {
	if s == nil {
		return nil
	}
	return &dtos.School{
		NPSN:           s.NPSN,
		NPSNVerifiedAt: s.NPSNVerifiedAt,
		Name:           s.Name,
		Level:          s.Level,
		Ownership:      s.Ownership,
		ProfilePicture: s.ProfilePicture,
		ProfileBanner:  s.ProfileBanner,
		Accreditation:  s.Accreditation,
		EstablishedAt:  s.EstablishedAt,
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
		CreatedAt:      s.CreatedAt,
	}
}
