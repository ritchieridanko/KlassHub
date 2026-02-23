package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type AuthRepository interface {
	GetByIdentifier(ctx context.Context, identifier string) (a *models.Auth, err *ce.Error)
}

type authRepository struct {
	database databases.AuthDatabase
}

func NewAuthRepository(db databases.AuthDatabase) AuthRepository {
	return &authRepository{database: db}
}

func (r *authRepository) GetByIdentifier(ctx context.Context, identifier string) (*models.Auth, *ce.Error) {
	return r.database.GetByIdentifier(ctx, identifier)
}
