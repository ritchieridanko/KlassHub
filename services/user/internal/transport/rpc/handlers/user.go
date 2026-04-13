package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	apis.UnimplementedUserServiceServer
	uu usecases.UserUsecase
}

func NewUserHandler(uu usecases.UserUsecase) *UserHandler {
	return &UserHandler{uu: uu}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *apis.CreateUserRequest) (*apis.CreateUserResponse, error) {
	u, err := h.uu.CreateUser(
		ctx,
		&models.CreateUserReq{
			AuthID:       req.GetAuthId(),
			SchoolID:     req.GetSchoolId(),
			SchoolUserID: req.SchoolUserId,
			Role:         req.GetRole(),
			Name:         req.GetName(),
			Birthplace:   req.GetBirthplace(),
			Birthdate:    *utils.ToTime(req.GetBirthdate()),
			Sex:          req.GetSex(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateUserResponse{
		User: h.toUserAdmin(u),
	}, nil
}

func (h *UserHandler) GetMe(ctx context.Context, req *emptypb.Empty) (*apis.UserGetMeResponse, error) {
	u, err := h.uu.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	return &apis.UserGetMeResponse{
		User: h.toUser(u),
	}, nil
}

func (h *UserHandler) toUser(u *models.User) *apis.User {
	if u == nil {
		return nil
	}
	return &apis.User{
		Id:             u.ID.String(),
		SchoolUserId:   u.SchoolUserID,
		Role:           utils.ToTitlecase(u.Role),
		Name:           u.Name,
		Nickname:       u.Nickname,
		Birthplace:     utils.ToTitlecase(u.Birthplace),
		Birthdate:      utils.ToTimestamp(&u.Birthdate),
		Sex:            utils.ToTitlecase(u.Sex),
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		ProfileBanner:  u.ProfileBanner,
	}
}

func (h *UserHandler) toUserAdmin(u *models.User) *apis.UserAdmin {
	if u == nil {
		return nil
	}

	var createdBy *string
	if u.CreatedBy != nil {
		id := u.CreatedBy.String()
		createdBy = &id
	}

	return &apis.UserAdmin{
		Id:             u.ID.String(),
		SchoolUserId:   u.SchoolUserID,
		Role:           utils.ToTitlecase(u.Role),
		Name:           u.Name,
		Birthplace:     utils.ToTitlecase(u.Birthplace),
		Birthdate:      utils.ToTimestamp(&u.Birthdate),
		Sex:            utils.ToTitlecase(u.Sex),
		Phone:          u.Phone,
		ProfilePicture: u.ProfilePicture,
		CreatedBy:      createdBy,
		CreatedAt:      utils.ToTimestamp(&u.CreatedAt),
		UpdatedAt:      utils.ToTimestamp(&u.UpdatedAt),
	}
}
