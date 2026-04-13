package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/models"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserClient interface {
	GetMe(ctx context.Context) (u *models.User, err *ce.Error)
}

type userClient struct {
	client apis.UserServiceClient
}

func NewUserClient(c apis.UserServiceClient) UserClient {
	return &userClient{client: c}
}

func (c *userClient) GetMe(ctx context.Context) (*models.User, *ce.Error) {
	resp, err := c.client.GetMe(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "user"),
		)
	}
	return c.toUser(resp.GetUser()), nil
}

func (c *userClient) toUser(u *apis.User) *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		ID:             utils.ToUUIDMust(u.GetId()),
		SchoolUserID:   u.SchoolUserId,
		Role:           u.GetRole(),
		Name:           u.GetName(),
		Nickname:       u.Nickname,
		Birthplace:     u.GetBirthplace(),
		Birthdate:      utils.ToTime(u.GetBirthdate()),
		Sex:            u.GetSex(),
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}
