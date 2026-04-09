package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/models"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type AuthClient interface {
	UpdateSchool(ctx context.Context, req *models.UpdateSchoolReq) (a *models.Auth, at *models.AuthToken, err *ce.Error)
}

type authClient struct {
	client apis.AuthServiceClient
}

func NewAuthClient(c apis.AuthServiceClient) AuthClient {
	return &authClient{client: c}
}

func (c *authClient) UpdateSchool(ctx context.Context, req *models.UpdateSchoolReq) (*models.Auth, *models.AuthToken, *ce.Error) {
	resp, err := c.client.UpdateSchool(
		ctx,
		&apis.UpdateSchoolRequest{
			SchoolId:     req.SchoolID,
			RefreshToken: req.RefreshToken,
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
