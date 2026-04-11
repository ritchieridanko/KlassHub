package databases

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/user/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/user/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserDatabase interface {
	Create(ctx context.Context, data *models.CreateUserData) (u *models.User, err *ce.Error)
	GetByAuthID(ctx context.Context, authID int64) (u *models.User, err *ce.Error)
}

type userDatabase struct {
	database *database.Database
}

func NewUserDatabase(db *database.Database) UserDatabase {
	return &userDatabase{database: db}
}

func (d *userDatabase) Create(ctx context.Context, data *models.CreateUserData) (*models.User, *ce.Error) {
	query := `
		INSERT INTO users (
			id, auth_id, school_id, school_user_id, role,
			name, birthplace, birthdate, sex, created_by
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10
		)
		RETURNING
			id, school_user_id, role, name, nickname,
			birthplace, birthdate, sex, phone,
			profile_picture, profile_banner,
			created_by, created_at, updated_at
	`

	var u models.User
	err := d.database.Query(
		ctx, query,
		data.ID, data.AuthID, data.SchoolID, data.SchoolUserID,
		data.Role, data.Name, data.Birthplace, data.Birthdate,
		data.Sex, data.CreatedBy,
	).Scan(
		&u.ID, &u.SchoolUserID, &u.Role, &u.Name, &u.Nickname,
		&u.Birthplace, &u.Birthdate, &u.Sex, &u.Phone,
		&u.ProfilePicture, &u.ProfileBanner, &u.CreatedBy,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create user: %w", err),
			logger.NewField("auth_id", data.AuthID),
			logger.NewField("school_id", data.SchoolID),
			logger.NewField("role", data.Role),
		)
	}

	return &u, nil
}

func (d *userDatabase) GetByAuthID(ctx context.Context, authID int64) (*models.User, *ce.Error) {
	query := `
		SELECT
			id, school_user_id, role, name, nickname,
			birthplace, birthdate, sex, phone,
			profile_picture, profile_banner
		FROM
			users
		WHERE
			auth_id = $1
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var u models.User
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(
		&u.ID, &u.SchoolUserID, &u.Role, &u.Name, &u.Nickname,
		&u.Birthplace, &u.Birthdate, &u.Sex, &u.Phone,
		&u.ProfilePicture, &u.ProfileBanner,
	)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to get user by auth id: %w", err)
		authIDField := logger.NewField("auth_id", authID)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeUserNotFound,
				ce.MsgUserNotFound,
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

	return &u, nil
}
