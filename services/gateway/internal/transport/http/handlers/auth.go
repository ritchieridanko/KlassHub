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
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/validator"
)

type AuthHandler struct {
	auc       clients.AuthClient
	validator *validator.Validator
	cookie    *cookie.Cookie
}

func NewAuthHandler(auc clients.AuthClient, v *validator.Validator, c *cookie.Cookie) *AuthHandler {
	return &AuthHandler{
		auc:       auc,
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

	a, at, err := h.auc.Login(
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

func (h *AuthHandler) Logout(ctx *gin.Context) {
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

	logoutErr := h.auc.Logout(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, false, false, false)...,
		),
		refreshToken,
	)
	if logoutErr != nil {
		logoutErr.Bind(ctx)
		return
	}

	h.cookie.Unset(
		ctx,
		constants.CookieKeyRefreshToken,
		"/",
	)

	utils.SetResponse[any](ctx, http.StatusNoContent, "", nil)
}

func (h *AuthHandler) CreateSchoolAuth(ctx *gin.Context) {
	var payload dtos.CreateSchoolAuthRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ce.NewError(ce.CodeInvalidPayload, ce.MsgInvalidPayload, err).Bind(ctx)
		return
	}

	a, at, err := h.auc.CreateSchoolAuth(
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

func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	var payload dtos.ChangePasswordRequest
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

	a, err := h.auc.ChangePassword(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, true)...,
		),
		&models.ChangePasswordReq{
			OldPassword: payload.OldPassword,
			NewPassword: payload.NewPassword,
		},
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Password changed successfully",
		dtos.ChangePasswordResponse{
			Auth: h.toAuth(a),
		},
	)
}

func (h *AuthHandler) ResendVerification(ctx *gin.Context) {
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

	email, err := h.auc.ResendVerification(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, false)...,
		),
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"Verification resent successfully",
		dtos.ResendVerificationResponse{
			Email: email,
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

	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		ce.NewError(ce.CodeRefreshTokenNotFound, ce.MsgInvalidSession, nil).Bind(ctx)
		return
	}

	a, at, verifyErr := h.auc.VerifyEmail(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.Auth(authCtx, true, true, true, false)...,
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

func (h *AuthHandler) RotateAuthToken(ctx *gin.Context) {
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

	at, rotateErr := h.auc.RotateAuthToken(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
		),
		refreshToken,
	)
	if rotateErr != nil {
		rotateErr.Bind(ctx)
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
		"Token refreshed successfully",
		dtos.RotateAuthTokenResponse{
			AccessToken: h.toAccessToken(at),
		},
	)
}

func (h *AuthHandler) IsEmailAvailable(ctx *gin.Context) {
	var params dtos.EmailAvailabilityCheckRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ce.NewError(ce.CodeInvalidParams, ce.MsgInvalidParams, err).Bind(ctx)
		return
	}

	available, err := h.auc.IsEmailAvailable(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
		),
		params.Email,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"OK",
		dtos.EmailAvailabilityCheckResponse{
			IsAvailable: available,
		},
	)
}

func (h *AuthHandler) IsUsernameAvailable(ctx *gin.Context) {
	var params dtos.UsernameAvailabilityCheckRequest
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

	available, err := h.auc.IsUsernameAvailable(
		metadata.ToOutgoingCtx(
			ctx.Request.Context(),
			metadata.NewPair(
				constants.MDKeyRole,
				authCtx.Role,
			),
		),
		params.Username,
	)
	if err != nil {
		err.Bind(ctx)
		return
	}

	utils.SetResponse(
		ctx,
		http.StatusOK,
		"OK",
		dtos.UsernameAvailabilityCheckResponse{
			IsAvailable: available,
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
