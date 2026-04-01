package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type SessionRepository interface {
	Create(ctx context.Context, data *models.CreateSessionData) (err *ce.Error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (s *models.Session, err *ce.Error)
	Revoke(ctx context.Context, params *models.RevokeSessionParams) (s *models.Session, err *ce.Error)
	RevokeActive(ctx context.Context, params *models.RevokeActiveSessionParams) (sessionID int64, err *ce.Error)
}

type sessionRepository struct {
	database databases.SessionDatabase
}

func NewSessionRepository(db databases.SessionDatabase) SessionRepository {
	return &sessionRepository{database: db}
}

func (r *sessionRepository) Create(ctx context.Context, data *models.CreateSessionData) *ce.Error {
	return r.database.Create(ctx, data)
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, *ce.Error) {
	return r.database.GetByRefreshToken(ctx, refreshToken)
}

func (r *sessionRepository) Revoke(ctx context.Context, params *models.RevokeSessionParams) (*models.Session, *ce.Error) {
	return r.database.Revoke(ctx, params)
}

func (r *sessionRepository) RevokeActive(ctx context.Context, params *models.RevokeActiveSessionParams) (int64, *ce.Error) {
	return r.database.RevokeActive(ctx, params)
}
