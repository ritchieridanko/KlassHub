package handlers

import (
	"errors"
	"net/http"
	"strconv"

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
	ac        clients.AuthClient
	validator *validator.Validator
	cookie    *cookie.Cookie
}

func NewAuthHandler(ac clients.AuthClient, v *validator.Validator, c *cookie.Cookie) *AuthHandler {
	return &AuthHandler{
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
		http.StatusOK,
		"Logged in successfully",
		dtos.LoginResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
		},
	)
}

func (h *AuthHandler) CreateSchoolAuth(ctx *gin.Context) {
	var payload dtos.CreateSchoolAuthRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	a, at, err := h.ac.CreateSchoolAuth(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
		),
		&models.CreateSchoolAuthReq{
			Email:    payload.Email,
			Password: payload.Password,
		},
	)
	if err != nil {
		err.Bind(ctx)
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
		"Registered successfully",
		dtos.CreateSchoolAuthResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
		},
	)
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	var params dtos.VerifyEmailRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
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

	a, at, verifyErr := h.ac.VerifyEmail(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.NewPair(
				constants.MDKeyAuthID,
				strconv.FormatInt(authCtx.AuthID, 10),
			),
			metadata.NewPair(
				constants.MDKeySchoolID,
				strconv.FormatInt(authCtx.SchoolID, 10),
			),
			metadata.NewPair(
				constants.MDKeyRole,
				authCtx.Role,
			),
		),
		&models.VerifyEmailReq{
			VerificationToken: params.VerificationToken,
			RefreshToken:      refreshToken,
		},
	)
	if verifyErr != nil {
		verifyErr.Bind(ctx)
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
		http.StatusOK,
		"Email verified successfully",
		dtos.VerifyEmailResponse{
			Auth:        h.toAuth(a),
			AccessToken: h.toAccessToken(at),
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
		SchoolExists:      a.SchoolExists,
		PasswordChangedAt: a.PasswordChangedAt,
	}
}

func (h *AuthHandler) toAccessToken(at *models.AuthToken) *dtos.AccessToken {
	if at == nil || at.AccessToken == nil {
		return nil
	}
	return &dtos.AccessToken{
		Token:     at.AccessToken.Token,
		ExpiresIn: at.AccessToken.ExpiresIn,
	}
}
