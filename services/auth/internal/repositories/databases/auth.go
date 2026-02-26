package databases

import (
	"context"
	"errors"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/models"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
)

type AuthDatabase interface {
	GetByIdentifier(ctx context.Context, identifier string) (a *models.Auth, err *ce.Error)
}

type authDatabase struct {
	database *database.Database
}

func NewAuthDatabase(db *database.Database) AuthDatabase {
	return &authDatabase{database: db}
}

func (d *authDatabase) GetByIdentifier(ctx context.Context, identifier string) (*models.Auth, *ce.Error) {
	query := `
		SELECT
			id, email, username, password, is_school,
			email_verified_at, password_changed_at
		FROM auth
		WHERE
			(
				email = $1
				OR username = $1
			)
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		identifier,
	).Scan(
		&a.ID, &a.Email, &a.Username, &a.Password,
		&a.IsSchool, &a.EmailVerifiedAt, &a.PasswordChangedAt,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(ce.CodeAuthNotFound, ce.MsgAuthNotFound, err)
		}
		return nil, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, err)
	}

	return &a, nil
}
