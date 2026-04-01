package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type TokenRepository interface {
	CreateVerification(ctx context.Context, data *models.CreateVerificationTokenData) (err *ce.Error)
	UseVerification(ctx context.Context, token string) (authID int64, err *ce.Error)
}

type tokenRepository struct {
	cache caches.TokenCache
}

func NewTokenRepository(cc caches.TokenCache) TokenRepository {
	return &tokenRepository{cache: cc}
}

func (r *tokenRepository) CreateVerification(ctx context.Context, data *models.CreateVerificationTokenData) *ce.Error {
	return r.cache.CreateVerification(ctx, data)
}

func (r *tokenRepository) UseVerification(ctx context.Context, token string) (int64, *ce.Error) {
	return r.cache.UseVerification(ctx, token)
}
