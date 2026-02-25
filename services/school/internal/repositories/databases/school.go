package databases

import (
	"context"
	"errors"

	"github.com/ritchieridanko/klasshub/services/school/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
)

type SchoolDatabase interface {
	GetID(ctx context.Context, authID int64) (schoolID int64, err *ce.Error)
}

type schoolDatabase struct {
	database *database.Database
}

func NewSchoolDatabase(db *database.Database) SchoolDatabase {
	return &schoolDatabase{database: db}
}

func (d *schoolDatabase) GetID(ctx context.Context, authID int64) (int64, *ce.Error) {
	query := "SELECT id FROM schools WHERE auth_id = $1 AND deleted_at IS NULL"
	if d.database.WithinTx(ctx) {
		query += " FOR UPDATE"
	}

	var schoolID int64
	err := d.database.Query(
		ctx, query,
		authID,
	).Scan(&schoolID)
	if err != nil {
		authIDField := logger.NewField("auth_id", authID)

		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, ce.NewError(ce.CodeSchoolNotFound, ce.MsgSchoolNotFound, err, authIDField)
		}
		return 0, ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, err, authIDField)
	}

	return schoolID, nil
}
