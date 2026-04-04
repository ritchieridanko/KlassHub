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

type AuthDatabase interface {
	Create(ctx context.Context, data *models.CreateAuthData) (a *models.Auth, err *ce.Error)
	GetByID(ctx context.Context, authID int64) (a *models.Auth, err *ce.Error)
	GetByIdentifier(ctx context.Context, identifier string) (a *models.Auth, err *ce.Error)
	UpdatePassword(ctx context.Context, authID int64, newPassword string) (a *models.Auth, err *ce.Error)
	SetVerified(ctx context.Context, authID int64) (a *models.Auth, err *ce.Error)
	IsEmailRegistered(ctx context.Context, email string) (exists bool, err *ce.Error)
}

type authDatabase struct {
	database *database.Database
}

func NewAuthDatabase(db *database.Database) AuthDatabase {
	return &authDatabase{database: db}
}

func (d *authDatabase) Create(ctx context.Context, data *models.CreateAuthData) (*models.Auth, *ce.Error) {
	query := `
		INSERT INTO auth (
			school_id, email, username,
			password, role, verified_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING
			id, school_id, email, username,
			role, verified_at, password_changed_at
	`

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		data.SchoolID, data.Email, data.Username,
		data.Password, data.Role, data.VerifiedAt,
	).Scan(
		&a.ID, &a.SchoolID, &a.Email, &a.Username,
		&a.Role, &a.VerifiedAt, &a.PasswordChangedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create auth: %w", err),
			logger.NewField("school_id", data.SchoolID),
			logger.NewField("role", data.Role),
		)
	}

	return &a, nil
}

func (d *authDatabase) GetByID(ctx context.Context, authID int64) (*models.Auth, *ce.Error) {
	query := `
		SELECT
			id, school_id, email, username, password,
			role, verified_at, password_changed_at
		FROM
			auth
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&a.ID, &a.SchoolID, &a.Email, &a.Username, &a.Password,
		&a.Role, &a.VerifiedAt, &a.PasswordChangedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get auth by id: %w", err)
		authIDField := logger.NewField("auth_id", authID)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAuthNotFound,
				ce.MsgAuthNotFound,
				wrappedErr,
				authIDField,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
			authIDField,
		)
	}

	return &a, nil
}

func (d *authDatabase) GetByIdentifier(ctx context.Context, identifier string) (*models.Auth, *ce.Error) {
	query := `
		SELECT
			id, school_id, email, username, password,
			role, verified_at, password_changed_at
		FROM
			auth
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
		&a.ID, &a.SchoolID, &a.Email, &a.Username, &a.Password,
		&a.Role, &a.VerifiedAt, &a.PasswordChangedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get auth by identifier: %w", err)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(ce.CodeAuthNotFound, ce.MsgAuthNotFound, wrappedErr)
		}
		return nil, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, wrappedErr)
	}

	return &a, nil
}

func (d *authDatabase) UpdatePassword(ctx context.Context, authID int64, newPassword string) (*models.Auth, *ce.Error) {
	query := `
		UPDATE auth
		SET
			password = $1,
			password_changed_at = NOW(),
			updated_at = NOW()
		WHERE
			id = $2
			AND deleted_at IS NULL
		RETURNING
			id, school_id, email, username,
			role, verified_at, password_changed_at
	`

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		newPassword, authID,
	).Scan(
		&a.ID, &a.SchoolID, &a.Email, &a.Username,
		&a.Role, &a.VerifiedAt, &a.PasswordChangedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update password: %w", err)
		authIDField := logger.NewField("auth_id", authID)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAuthNotFound,
				ce.MsgAuthNotFound,
				wrappedErr,
				authIDField,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
			authIDField,
		)
	}

	return &a, nil
}

func (d *authDatabase) SetVerified(ctx context.Context, authID int64) (*models.Auth, *ce.Error) {
	query := `
		UPDATE auth
		SET
			verified_at = NOW(),
			updated_at = NOW()
		WHERE
			id = $1
			AND deleted_at IS NULL
		RETURNING
			id, school_id, email, username,
			role, verified_at, password_changed_at
	`

	var a models.Auth
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&a.ID, &a.SchoolID, &a.Email, &a.Username,
		&a.Role, &a.VerifiedAt, &a.PasswordChangedAt,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to set auth verified: %w", err)
		authIDField := logger.NewField("auth_id", authID)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeAuthNotFound,
				ce.MsgAuthNotFound,
				wrappedErr,
				authIDField,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			wrappedErr,
			authIDField,
		)
	}

	return &a, nil
}

func (d *authDatabase) IsEmailRegistered(ctx context.Context, email string) (bool, *ce.Error) {
	query := "SELECT 1 FROM auth WHERE email = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var exists int
	err := d.database.Query(
		ctx, query,
		email,
	).Scan(
		&exists,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return false, nil
		}
		return false, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to check if email is registered: %w", err),
		)
	}

	return true, nil
}
