package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type UserHandler struct {
	apis.UnimplementedUserServiceServer
	uu usecases.UserUsecase
}

func NewUserHandler(uu usecases.UserUsecase) *UserHandler {
	return &UserHandler{uu: uu}
}

func (h *UserHandler) GetUser(ctx context.Context, req *apis.GetUserRequest) (*apis.GetUserResponse, error) {
	u, err := h.uu.GetUser(
		ctx,
		&models.GetUserRequest{
			AuthID:   req.GetAuthId(),
			SchoolID: req.GetSchoolId(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.GetUserResponse{User: h.toUser(u)}, nil
}

func (h *UserHandler) GetSchoolAndRole(ctx context.Context, req *apis.GetSchoolAndRoleRequest) (*apis.GetSchoolAndRoleResponse, error) {
	schoolID, role, err := h.uu.GetSchoolAndRole(ctx, req.GetAuthId())
	if err != nil {
		return nil, err
	}
	return &apis.GetSchoolAndRoleResponse{
		SchoolId: schoolID,
		Role:     role,
	}, nil
}

func (h *UserHandler) toUser(u *models.User) *apis.User {
	if u == nil {
		return nil
	}

	var createdBy *string
	if u.CreatedBy != nil {
		v := u.CreatedBy.String()
		createdBy = &v
	}

	return &apis.User{
		Id:             u.ID.String(),
		SchoolUserId:   utils.WrapString(u.SchoolUserID),
		Role:           u.Role,
		Name:           u.Name,
		Nickname:       utils.WrapString(u.Nickname),
		Birthplace:     utils.WrapString(u.Birthplace),
		Birthdate:      utils.WrapTime(u.Birthdate),
		Sex:            utils.WrapString(u.Sex),
		Phone:          utils.WrapString(u.Phone),
		ProfilePicture: utils.WrapString(u.ProfilePicture),
		ProfileBanner:  utils.WrapString(u.ProfileBanner),
		CreatedBy:      utils.WrapString(createdBy),
		CreatedByName:  utils.WrapString(u.CreatedByName),
		CreatedAt:      utils.WrapTime(u.CreatedAt),
		UpdatedAt:      utils.WrapTime(u.UpdatedAt),
		DeletedAt:      utils.WrapTime(u.DeletedAt),
	}
}
