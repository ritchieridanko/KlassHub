package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AuthClient interface {
	Login(ctx context.Context, req *models.LoginReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	Logout(ctx context.Context, refreshToken string) (err *ce.Error)
	CreateSchoolAuth(ctx context.Context, req *models.CreateSchoolAuthReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authClient struct {
	client apis.AuthServiceClient
}

func NewAuthClient(c apis.AuthServiceClient) AuthClient {
	return &authClient{client: c}
}

func (c *authClient) Login(ctx context.Context, req *models.LoginReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.Login(
		ctx,
		&apis.LoginRequest{
			Identifier: req.Identifier,
			Password:   req.Password,
		},
	)
	if err != nil {
		return nil, nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "auth"),
		)
	}
	return c.toAuth(resp.GetAuth()), c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *authClient) Logout(ctx context.Context, refreshToken string) *ce.Error {
	_, err := c.client.Logout(
		ctx,
		&apis.LogoutRequest{
			RefreshToken: refreshToken,
		},
	)
	if err != nil {
		return ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "auth"),
		)
	}
	return nil
}

func (c *authClient) CreateSchoolAuth(ctx context.Context, req *models.CreateSchoolAuthReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.CreateSchoolAuth(
		ctx,
		&apis.CreateSchoolAuthRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "auth"),
		)
	}
	return c.toAuth(resp.GetAuth()), c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *authClient) VerifyEmail(ctx context.Context, req *models.VerifyEmailReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.VerifyEmail(
		ctx,
		&apis.VerifyEmailRequest{
			VerificationToken: req.VerificationToken,
			RefreshToken:      req.RefreshToken,
		},
	)
	if err != nil {
		return nil, nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "auth"),
		)
	}
	return c.toAuth(resp.GetAuth()), c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *authClient) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	resp, err := c.client.IsEmailAvailable(
		ctx,
		&apis.EmailAvailabilityCheckRequest{
			Email: email,
		},
	)
	if err != nil {
		return false, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "auth"),
		)
	}
	return resp.GetIsAvailable(), nil
}

func (c *authClient) toAuth(a *apis.Auth) *models.Auth {
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

func (c *authClient) toAuthToken(at *apis.AuthToken) *models.AuthToken {
	if at == nil {
		return nil
	}
	return &models.AuthToken{
		AccessToken:  c.toAccessToken(at.GetAccessToken()),
		RefreshToken: c.toRefreshToken(at.GetRefreshToken()),
	}
}

func (c *authClient) toAccessToken(at *apis.AccessToken) *models.AccessToken {
	if at == nil {
		return nil
	}
	return &models.AccessToken{
		Token:     at.GetToken(),
		ExpiresIn: at.GetExpiresIn(),
	}
}

func (c *authClient) toRefreshToken(rt *apis.RefreshToken) *models.RefreshToken {
	if rt == nil {
		return nil
	}
	return &models.RefreshToken{
		Token:     rt.GetToken(),
		ExpiresIn: rt.GetExpiresIn(),
	}
}
