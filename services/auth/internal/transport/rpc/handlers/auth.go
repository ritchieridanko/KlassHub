package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AuthHandler struct {
	apis.UnimplementedAuthServiceServer
	au usecases.AuthUsecase
}

func NewAuthHandler(au usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{au: au}
}

func (h *AuthHandler) Login(ctx context.Context, req *apis.LoginRequest) (*apis.LoginResponse, error) {
	a, at, err := h.au.Login(
		ctx,
		&models.LoginReq{
			Identifier: req.GetIdentifier(),
			Password:   req.GetPassword(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.LoginResponse{
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) toAuth(a *models.Auth) *apis.Auth {
	if a == nil {
		return nil
	}
	return &apis.Auth{
		Email:             a.Email,
		Username:          a.Username,
		Role:              a.Role,
		IsVerified:        a.IsVerified(),
		PasswordChangedAt: utils.ToTimestamp(a.PasswordChangedAt),
	}
}

func (h *AuthHandler) toAuthToken(at *models.AuthToken) *apis.AuthToken {
	if at == nil {
		return nil
	}
	return &apis.AuthToken{
		AccessToken:  h.toAccessToken(at.AccessToken),
		RefreshToken: h.toRefreshToken(at.RefreshToken),
	}
}

func (h *AuthHandler) toAccessToken(at *models.AccessToken) *apis.AccessToken {
	if at == nil {
		return nil
	}
	return &apis.AccessToken{
		Token:     at.Token,
		ExpiresIn: at.ExpiresIn,
	}
}

func (h *AuthHandler) toRefreshToken(rt *models.RefreshToken) *apis.RefreshToken {
	if rt == nil {
		return nil
	}
	return &apis.RefreshToken{
		Token:     rt.Token,
		ExpiresIn: rt.ExpiresIn,
	}
}
