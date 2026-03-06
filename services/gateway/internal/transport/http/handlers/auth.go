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
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/metadata"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/validator"
)

type AuthHandler struct {
	hostname  string
	tld       string
	ac        clients.AuthClient
	validator *validator.Validator
	cookie    *cookie.Cookie
}

func NewAuthHandler(hostname, tld string, ac clients.AuthClient, v *validator.Validator, c *cookie.Cookie) *AuthHandler {
	return &AuthHandler{
		hostname:  hostname,
		tld:       tld,
		ac:        ac,
		validator: v,
		cookie:    c,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var payload dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	ua, ip := ctx.Request.UserAgent(), ctx.ClientIP()
	if ok, why := h.validator.UserAgent(ua); !ok {
		ce.NewError(ce.CodeInvalidRequestMetadata, why, nil).Bind(ctx)
		return
	}
	if ok, why := h.validator.IPAddress(ip); !ok {
		ce.NewError(ce.CodeInvalidRequestMetadata, why, nil).Bind(ctx)
		return
	}

	a, at, err := h.ac.Login(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.NewPair(
				constants.MDKeyUserAgent,
				ua,
			),
			metadata.NewPair(
				constants.MDKeyIPAddress,
				ip,
			),
			metadata.NewPair(
				constants.MDKeySubdomain,
				utils.CtxSubdomain(ctx.Request.Context()),
			),
		),
		&models.LoginReq{
			Identifier: payload.Identifier,
			Password:   payload.Password,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	h.cookie.Set(
		ctx,
		constants.CookieKeyRefreshToken,
		at.RefreshToken.Token,
		"/",
		int(at.RefreshToken.ExpiresIn),
	)
	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Logged in successfully",
		dtos.LoginResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at.AccessToken),
		},
	)
}

func (h *AuthHandler) toAuth(a *models.Auth) *dtos.Auth {
	if a == nil {
		return nil
	}
	return &dtos.Auth{
		Email:             a.Email,
		Username:          a.Username,
		Role:              a.Role,
		IsVerified:        a.IsVerified,
		PasswordChangedAt: a.PasswordChangedAt,
	}
}

func (h *AuthHandler) toAccessToken(at *models.AccessToken) *dtos.AccessToken {
	if at == nil {
		return nil
	}
	return &dtos.AccessToken{
		Token:     at.Token,
		ExpiresIn: at.ExpiresIn,
	}
}
