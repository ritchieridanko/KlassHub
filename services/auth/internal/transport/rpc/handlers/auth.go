package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AuthHandler struct {
	apis.UnimplementedAuthServiceServer
	au     usecases.AuthUsecase
	logger *logger.Logger
}

func NewAuthHandler(au usecases.AuthUsecase, l *logger.Logger) *AuthHandler {
	return &AuthHandler{au: au, logger: l}
}

func (h *AuthHandler) Login(ctx context.Context, req *apis.LoginRequest) (*apis.LoginResponse, error) {
	ua, ip := utils.CtxRequestMeta(ctx)
	if ua == "" || ip == "" {
		h.logger.Warn(
			ctx,
			"incomplete request meta",
			logger.NewField("user_agent", ua),
			logger.NewField("ip_address", ip),
			logger.NewField("error_code", ce.CodeInvalidRequestMeta),
		)
	}

	a, at, err := h.au.Login(
		ctx,
		&models.LoginRequest{
			Identifier: req.GetIdentifier(),
			Password:   req.GetPassword(),
			Subdomain:  req.GetSubdomain(),
			RequestMeta: models.RequestMeta{
				UserAgent: ua,
				IPAddress: ip,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &apis.LoginResponse{
		AuthToken: h.toAuthToken(at),
		Auth:      h.toAuth(a),
	}, nil
}

func (h *AuthHandler) toAuthToken(at *models.AuthToken) *apis.AuthToken {
	if at == nil {
		return nil
	}
	return &apis.AuthToken{
		AccessToken:           at.AccessToken,
		RefreshToken:          at.RefreshToken,
		AccessTokenExpiresIn:  at.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: at.RefreshTokenExpiresIn,
	}
}

func (h *AuthHandler) toAuth(a *models.Auth) *apis.Auth {
	if a == nil {
		return nil
	}
	return &apis.Auth{
		Role:              utils.WrapString(&a.Role),
		Email:             utils.WrapString(a.Email),
		Username:          utils.WrapString(a.Username),
		EmailVerifiedAt:   utils.WrapTime(a.EmailVerifiedAt),
		PasswordChangedAt: utils.WrapTime(a.PasswordChangedAt),
	}
}
