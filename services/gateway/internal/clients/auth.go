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
	Login(ctx context.Context, req *models.LoginRequest) (a *models.Auth, at *models.AuthToken, err *ce.Error)
}

type authClient struct {
	client apis.AuthServiceClient
}

func NewAuthClient(c apis.AuthServiceClient) AuthClient {
	return &authClient{client: c}
}

func (c *authClient) Login(ctx context.Context, req *models.LoginRequest) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.Login(
		ctx,
		&apis.LoginRequest{
			Identifier: req.Identifier,
			Password:   req.Password,
			Subdomain:  req.Subdomain,
		},
	)
	if err != nil {
		return nil, nil, ce.FromGRPCErr(
			err,
		).AppendFields(
			logger.NewField("service", "auth"),
			logger.NewField("subdomain", req.Subdomain),
		)
	}
	return c.toAuth(resp.GetAuth()), c.toAuthToken(resp.GetAuthToken()), nil
}

func (c *authClient) toAuth(a *apis.Auth) *models.Auth {
	if a == nil {
		return nil
	}
	return &models.Auth{
		Role:              utils.UnwrapString(a.GetRole()),
		Email:             utils.UnwrapString(a.GetEmail()),
		Username:          utils.UnwrapString(a.GetUsername()),
		EmailVerifiedAt:   utils.UnwrapTimestamp(a.GetEmailVerifiedAt()),
		PasswordChangedAt: utils.UnwrapTimestamp(a.GetPasswordChangedAt()),
	}
}

func (c *authClient) toAuthToken(at *apis.AuthToken) *models.AuthToken {
	if at == nil {
		return nil
	}
	return &models.AuthToken{
		AccessToken:           at.GetAccessToken(),
		RefreshToken:          at.GetRefreshToken(),
		AccessTokenExpiresIn:  at.GetAccessTokenExpiresIn(),
		RefreshTokenExpiresIn: at.GetRefreshTokenExpiresIn(),
	}
}
