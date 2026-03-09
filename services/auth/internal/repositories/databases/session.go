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
	RevokeActive(ctx context.Context, params *models.RevokeSessionParams) (sessionID int64, err *ce.Error)
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

func (d *sessionDatabase) RevokeActive(ctx context.Context, params *models.RevokeSessionParams) (int64, *ce.Error) {
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
