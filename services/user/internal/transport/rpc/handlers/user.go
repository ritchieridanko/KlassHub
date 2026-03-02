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

func (h *UserHandler) GetUserAuthInfo(ctx context.Context, req *apis.GetUserAuthInfoRequest) (*apis.GetUserAuthInfoResponse, error) {
	uai, err := h.uu.GetUserAuthInfo(
		ctx,
		&models.GetUserAuthInfoRequest{
			AuthID: req.GetAuthId(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.GetUserAuthInfoResponse{AuthInfo: h.toUserAuthInfo(uai)}, nil
}

func (h *UserHandler) toUser(u *models.User) *apis.User {
	if u == nil {
		return nil
	}
	return &apis.User{
		Id:             u.ID.String(),
		SchoolUserId:   u.SchoolUserID,
		Role:           u.Role,
		Name:           u.Name,
		Nickname:       u.Nickname,
		Birthplace:     u.Birthplace,
		Birthdate:      utils.ToTimestamp(&u.Birthdate),
		Sex:            u.Sex,
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}

func (h *UserHandler) toUserAuthInfo(uai *models.UserAuthInfo) *apis.UserAuthInfo {
	if uai == nil {
		return nil
	}
	return &apis.UserAuthInfo{
		SchoolId: uai.SchoolID,
		Role:     uai.Role,
	}
}
