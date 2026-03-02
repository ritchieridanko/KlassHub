package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserUsecase interface {
	GetUser(ctx context.Context, req *models.GetUserRequest) (u *models.User, err *ce.Error)
	GetUserAuthInfo(ctx context.Context, req *models.GetUserAuthInfoRequest) (uai *models.UserAuthInfo, err *ce.Error)
}

type userUsecase struct {
	ur repositories.UserRepository
}

func NewUserUsecase(ur repositories.UserRepository) UserUsecase {
	return &userUsecase{ur: ur}
}

func (u *userUsecase) GetUser(ctx context.Context, req *models.GetUserRequest) (*models.User, *ce.Error) {
	return u.ur.Get(
		ctx,
		&models.GetUser{
			AuthID:   req.AuthID,
			SchoolID: req.SchoolID,
		},
	)
}

func (u *userUsecase) GetUserAuthInfo(ctx context.Context, req *models.GetUserAuthInfoRequest) (*models.UserAuthInfo, *ce.Error) {
	return u.ur.GetAuthInfo(ctx, &models.GetUserAuthInfo{AuthID: req.AuthID})
}
