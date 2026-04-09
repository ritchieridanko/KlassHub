package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (h *AuthHandler) Logout(ctx context.Context, req *apis.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.au.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) CreateSchoolAuth(ctx context.Context, req *apis.CreateSchoolAuthRequest) (*apis.CreateSchoolAuthResponse, error) {
	a, at, err := h.au.CreateSchoolAuth(
		ctx,
		&models.CreateSchoolAuthReq{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateSchoolAuthResponse{
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) UpdateSchool(ctx context.Context, req *apis.UpdateSchoolRequest) (*apis.UpdateSchoolResponse, error) {
	a, at, err := h.au.UpdateSchool(
		ctx,
		&models.UpdateSchoolReq{
			SchoolID:     req.GetSchoolId(),
			RefreshToken: req.GetRefreshToken(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.UpdateSchoolResponse{
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *apis.ChangePasswordRequest) (*apis.ChangePasswordResponse, error) {
	a, err := h.au.ChangePassword(
		ctx,
		&models.ChangePasswordReq{
			OldPassword: req.GetOldPassword(),
			NewPassword: req.GetNewPassword(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.ChangePasswordResponse{
		Auth: h.toAuth(a),
	}, nil
}

func (h *AuthHandler) ResendVerification(ctx context.Context, req *emptypb.Empty) (*apis.ResendVerificationResponse, error) {
	email, err := h.au.ResendVerification(ctx)
	if err != nil {
		return nil, err
	}
	return &apis.ResendVerificationResponse{
		Email: email,
	}, nil
}

func (h *AuthHandler) VerifyEmail(ctx context.Context, req *apis.VerifyEmailRequest) (*apis.VerifyEmailResponse, error) {
	a, at, err := h.au.VerifyEmail(
		ctx,
		&models.VerifyEmailReq{
			VerificationToken: req.GetVerificationToken(),
			RefreshToken:      req.GetRefreshToken(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.VerifyEmailResponse{
		Auth:      h.toAuth(a),
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) RotateAuthToken(ctx context.Context, req *apis.RotateAuthTokenRequest) (*apis.RotateAuthTokenResponse, error) {
	at, err := h.au.RotateAuthToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return &apis.RotateAuthTokenResponse{
		AuthToken: h.toAuthToken(at),
	}, nil
}

func (h *AuthHandler) IsEmailAvailable(ctx context.Context, req *apis.EmailAvailabilityCheckRequest) (*apis.EmailAvailabilityCheckResponse, error) {
	available, err := h.au.IsEmailAvailable(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &apis.EmailAvailabilityCheckResponse{
		IsAvailable: available,
	}, nil
}

func (h *AuthHandler) IsUsernameAvailable(ctx context.Context, req *apis.UsernameAvailabilityCheckRequest) (*apis.UsernameAvailabilityCheckResponse, error) {
	available, err := h.au.IsUsernameAvailable(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}
	return &apis.UsernameAvailabilityCheckResponse{
		IsAvailable: available,
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
		SchoolExists:      a.SchoolExists(),
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
