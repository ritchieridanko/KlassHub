package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type AuthRepository interface {
	Create(ctx context.Context, data *models.CreateAuthData) (a *models.Auth, err *ce.Error)
	GetByID(ctx context.Context, authID int64) (a *models.Auth, err *ce.Error)
	GetByIdentifier(ctx context.Context, identifier string) (a *models.Auth, err *ce.Error)
	UpdatePassword(ctx context.Context, authID int64, newPassword string) (a *models.Auth, err *ce.Error)
	SetVerified(ctx context.Context, authID int64) (a *models.Auth, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authRepository struct {
	database databases.AuthDatabase
	cache    caches.AuthCache
}

func NewAuthRepository(db databases.AuthDatabase, cc caches.AuthCache) AuthRepository {
	return &authRepository{database: db, cache: cc}
}

func (r *authRepository) Create(ctx context.Context, data *models.CreateAuthData) (*models.Auth, *ce.Error) {
	return r.database.Create(ctx, data)
}

func (r *authRepository) GetByID(ctx context.Context, authID int64) (*models.Auth, *ce.Error) {
	return r.database.GetByID(ctx, authID)
}

func (r *authRepository) GetByIdentifier(ctx context.Context, identifier string) (*models.Auth, *ce.Error) {
	return r.database.GetByIdentifier(ctx, identifier)
}

func (r *authRepository) UpdatePassword(ctx context.Context, authID int64, newPassword string) (*models.Auth, *ce.Error) {
	return r.database.UpdatePassword(ctx, authID, newPassword)
}

func (r *authRepository) SetVerified(ctx context.Context, authID int64) (*models.Auth, *ce.Error) {
	return r.database.SetVerified(ctx, authID)
}

func (r *authRepository) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	registered, err := r.database.IsEmailRegistered(ctx, email)
	if err != nil {
		return false, err
	}
	if registered {
		return false, nil
	}

	reserved, err := r.cache.IsEmailReserved(ctx, email)
	if err != nil {
		return false, err
	}
	return !reserved, nil
}
