package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/clients"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/constants"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/dtos"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/cookie"
)

type AuthHandler struct {
	hostname string
	tld      string
	ac       clients.AuthClient
	cookie   *cookie.Cookie
}

func NewAuthHandler(hostname, tld string, ac clients.AuthClient, c *cookie.Cookie) *AuthHandler {
	return &AuthHandler{hostname: hostname, tld: tld, ac: ac, cookie: c}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var payload dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	subdomain, err := utils.CtxSubdomain(ctx, h.hostname, h.tld)
	if err != nil {
		ce.NewError(ce.CodeInvalidSubdomain, ce.MsgInvalidSubdomain, err).Bind(ctx)
		return
	}

	a, at, se := h.ac.Login(
		utils.ToOutgoingCtx(ctx, true),
		&models.LoginRequest{
			Identifier: payload.Identifier,
			Password:   payload.Password,
			Subdomain:  subdomain,
		},
	)
	if se != nil {
		se.Bind(ctx)
		return
	}

	h.cookie.Set(
		ctx,
		constants.CookieKeyRefreshToken,
		at.RefreshToken,
		"/",
		int(at.RefreshTokenExpiresIn),
	)
	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Logged in successfully",
		dtos.LoginResponse{
			Auth:      h.toAuth(a),
			AuthToken: h.toAuthToken(at),
		},
	)
}

func (h *AuthHandler) toAuth(a *models.Auth) *dtos.Auth {
	if a == nil {
		return nil
	}
	return &dtos.Auth{
		Role:              a.Role,
		Email:             a.Email,
		Username:          a.Username,
		EmailVerifiedAt:   a.EmailVerifiedAt,
		PasswordChangedAt: a.PasswordChangedAt,
	}
}

func (h *AuthHandler) toAuthToken(at *models.AuthToken) *dtos.AuthToken {
	if at == nil {
		return nil
	}
	return &dtos.AuthToken{
		AccessToken: at.AccessToken,
		ExpiresIn:   at.AccessTokenExpiresIn,
	}
}
