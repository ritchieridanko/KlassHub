package usecases

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserUsecase interface {
	GetUser(ctx context.Context, req *models.GetUserRequest) (u *models.User, err *ce.Error)
	GetSchoolAndRole(ctx context.Context, authID int64) (schoolID int64, role string, err *ce.Error)
}

type userUsecase struct {
	ur repositories.UserRepository
}

func NewUserUsecase(ur repositories.UserRepository) UserUsecase {
	return &userUsecase{ur: ur}
}

func (u *userUsecase) GetUser(ctx context.Context, req *models.GetUserRequest) (*models.User, *ce.Error) {
	return u.ur.GetByAuthID(
		ctx,
		&models.GetUser{
			AuthID:   req.AuthID,
			SchoolID: req.SchoolID,
		},
	)
}

func (u *userUsecase) GetSchoolAndRole(ctx context.Context, authID int64) (int64, string, *ce.Error) {
	return u.ur.GetSchoolAndRole(ctx, authID)
}
