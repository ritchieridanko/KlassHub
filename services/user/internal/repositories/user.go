package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserRepository interface {
	GetByAuthID(ctx context.Context, data *models.GetUser) (u *models.User, err *ce.Error)
}

type userRepository struct {
	database databases.UserDatabase
}

func NewUserRepository(db databases.UserDatabase) UserRepository {
	return &userRepository{database: db}
}

func (r *userRepository) GetByAuthID(ctx context.Context, data *models.GetUser) (*models.User, *ce.Error) {
	return r.database.GetByAuthID(ctx, data)
}
