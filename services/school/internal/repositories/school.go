package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
)

type SchoolRepository interface {
	GetID(ctx context.Context, params *models.GetSchoolID) (schoolID int64, err *ce.Error)
}

type schoolRepository struct {
	database databases.SchoolDatabase
}

func NewSchoolRepository(db databases.SchoolDatabase) SchoolRepository {
	return &schoolRepository{database: db}
}

func (r *schoolRepository) GetID(ctx context.Context, params *models.GetSchoolID) (int64, *ce.Error) {
	return r.database.GetID(ctx, params)
}
