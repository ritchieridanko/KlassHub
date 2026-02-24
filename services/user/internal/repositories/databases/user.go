package databases

import (
	"context"
	"errors"

	"github.com/ritchieridanko/klasshub/services/user/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/user/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/user/internal/models"
	"github.com/ritchieridanko/klasshub/services/user/internal/utils/ce"
)

type UserDatabase interface {
	GetByAuthID(ctx context.Context, data *models.GetUser) (u *models.User, err *ce.Error)
}

type userDatabase struct {
	database *database.Database
}

func NewUserDatabase(db *database.Database) UserDatabase {
	return &userDatabase{database: db}
}

func (d *userDatabase) GetByAuthID(ctx context.Context, data *models.GetUser) (*models.User, *ce.Error) {
	query := `
		SELECT
			id, school_user_id, role, name, nickname, birthplace,
			birthdate, sex, phone, profile_picture, profile_banner
		FROM users
		WHERE
			auth_id = $1
			AND school_id = $2
			AND deleted_at IS NULL
	`
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var u models.User
	err := d.database.Query(
		ctx, query,
		data.AuthID, data.SchoolID,
	).Scan(
		&u.ID, &u.SchoolUserID, &u.Role, &u.Name,
		&u.Nickname, &u.Birthplace, &u.Birthdate,
		&u.Sex, &u.Phone, &u.ProfilePicture, &u.ProfileBanner,
	)
	if err != nil {
		authIDField := logger.NewField("auth_id", data.AuthID)
		schoolIDField := logger.NewField("school_id", data.SchoolID)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeUserNotFound,
				ce.MsgResourceNotFound,
				err,
				authIDField,
				schoolIDField,
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			err,
			authIDField,
			schoolIDField,
		)
	}

	return &u, nil
}
