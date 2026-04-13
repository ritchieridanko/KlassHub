package clients

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/models"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type UserClient interface {
	CreateUser(ctx context.Context, req *models.CreateUserReq) (u *models.User, err *ce.Error)
}

type userClient struct {
	client apis.UserServiceClient
}

func NewUserClient(c apis.UserServiceClient) UserClient {
	return &userClient{client: c}
}

func (c *userClient) CreateUser(ctx context.Context, req *models.CreateUserReq) (*models.User, *ce.Error) {
	resp, err := c.client.CreateUser(
		ctx,
		&apis.CreateUserRequest{
			AuthId:       req.AuthID,
			SchoolId:     req.SchoolID,
			SchoolUserId: req.SchoolUserID,
			Role:         req.Role,
			Name:         req.Name,
			Birthplace:   req.Birthplace,
			Birthdate:    utils.ToTimestamp(req.Birthdate),
			Sex:          req.Sex,
		},
	)
	if err != nil {
		return nil, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "user"),
		)
	}
	return c.toUserFromUA(resp.GetUser()), nil
}

func (c *userClient) toUserFromUA(ua *apis.UserAdmin) *models.User {
	if ua == nil {
		return nil
	}

	var createdBy *uuid.UUID
	if ua.CreatedBy != nil {
		creator := utils.ToUUIDMust(ua.GetCreatedBy())
		createdBy = &creator
	}

	return &models.User{
		ID:             utils.ToUUIDMust(ua.GetId()),
		SchoolUserID:   ua.SchoolUserId,
		Role:           ua.GetRole(),
		Name:           ua.GetName(),
		Birthplace:     ua.GetBirthplace(),
		Birthdate:      utils.ToTime(ua.GetBirthdate()),
		Sex:            ua.GetSex(),
		Phone:          ua.Phone,
		ProfilePicture: ua.ProfilePicture,
		CreatedBy:      createdBy,
		CreatedAt:      utils.ToTime(ua.GetCreatedAt()),
		UpdatedAt:      utils.ToTime(ua.GetUpdatedAt()),
	}
}
