package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserRepository interface {
	Create(ctx context.Context, data *models.CreateUserData) (u *models.User, err *ce.Error)
	GetByAuthID(ctx context.Context, authID int64) (u *models.User, err *ce.Error)
}

type userRepository struct {
	database databases.UserDatabase
}

func NewUserRepository(db databases.UserDatabase) UserRepository {
	return &userRepository{database: db}
}

func (r *userRepository) Create(ctx context.Context, data *models.CreateUserData) (*models.User, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *userRepository) GetByAuthID(ctx context.Context, authID int64) (*models.User, *ce.Error) {
	return r.database.GetByAuthID(ctx, authID)
}
