package databases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type SessionDatabase interface {
	Create(ctx context.Context, data *models.CreateSessionData) (err *ce.Error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (s *models.Session, err *ce.Error)
	Revoke(ctx context.Context, params *models.RevokeSessionParams) (s *models.Session, err *ce.Error)
	RevokeActive(ctx context.Context, params *models.RevokeActiveSessionParams) (sessionID int64, err *ce.Error)
}

type sessionDatabase struct {
	database *database.Database
}

func NewSessionDatabase(db *database.Database) SessionDatabase {
	return &sessionDatabase{database: db}
}

func (d *sessionDatabase) Create(ctx context.Context, data *models.CreateSessionData) *ce.Error {
	query := `
		INSERT INTO sessions (
			parent_id, auth_id, refresh_token,
			user_agent, ip_address, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	err := d.database.Execute(
		ctx, query,
		data.ParentID, data.AuthID, data.RefreshToken,
		data.UserAgent, data.IPAddress, data.ExpiresAt,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create session: %w", err),
			logger.NewField("auth_id", data.AuthID),
		)
	}

	return nil
}

func (d *sessionDatabase) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, *ce.Error) {
	query := `
		SELECT
			id, auth_id, refresh_token,
			user_agent, ip_address, expires_at
		FROM
			sessions
		WHERE
			refresh_token = $1
			AND revoked_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var s models.Session
	err := d.database.Query(
		ctx, query,
		refreshToken,
	).Scan(
		&s.ID, &s.AuthID, &s.RefreshToken,
		&s.UserAgent, &s.IPAddress, &s.ExpiresAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get session by refresh token: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(ce.CodeSessionNotFound, ce.MsgSessionNotFound, wrappedErr)
		}
		return nil, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, wrappedErr)
	}

	return &s, nil
}

func (d *sessionDatabase) Revoke(ctx context.Context, params *models.RevokeSessionParams) (*models.Session, *ce.Error) {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE
			refresh_token = $1
			AND expires_at >= $2
			AND revoked_at IS NULL
		RETURNING
			id, auth_id, refresh_token,
			user_agent, ip_address
	`

	var s models.Session
	err := d.database.Query(
		ctx, query,
		params.RefreshToken, params.ExpiresAt,
	).Scan(
		&s.ID, &s.AuthID, &s.RefreshToken,
		&s.UserAgent, &s.IPAddress,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to revoke session: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeSessionNotFound,
				ce.MsgSessionNotFound,
				wrappedErr,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
		)
	}

	return &s, nil
}

func (d *sessionDatabase) RevokeActive(ctx context.Context, params *models.RevokeActiveSessionParams) (int64, *ce.Error) {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE
			auth_id = $1
			AND user_agent = $2
			AND ip_address = $3
			AND expires_at >= $4
			AND revoked_at IS NULL
		RETURNING id
	`

	var sessionID int64
	err := d.database.Query(
		ctx, query,
		params.AuthID, params.UserAgent,
		params.IPAddress, params.ExpiresAt,
	).Scan(
		&sessionID,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		return 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to revoke active session: %w", err),
			logger.NewField("auth_id", params.AuthID),
		)
	}

	return sessionID, nil
}
