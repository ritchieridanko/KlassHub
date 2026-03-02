package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserRepository interface {
	Get(ctx context.Context, params *models.GetUser) (u *models.User, err *ce.Error)
	GetAuthInfo(ctx context.Context, params *models.GetUserAuthInfo) (uai *models.UserAuthInfo, err *ce.Error)
}

type userRepository struct {
	database databases.UserDatabase
}

func NewUserRepository(db databases.UserDatabase) UserRepository {
	return &userRepository{database: db}
}

func (r *userRepository) Get(ctx context.Context, params *models.GetUser) (*models.User, *ce.Error) {
	return r.database.Get(ctx, params)
}

func (r *userRepository) GetAuthInfo(ctx context.Context, params *models.GetUserAuthInfo) (*models.UserAuthInfo, *ce.Error) {
	return r.database.GetAuthInfo(ctx, params)
}
