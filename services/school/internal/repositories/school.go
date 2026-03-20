package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/ce"
)

type SchoolRepository interface {
	Create(ctx context.Context, data *models.CreateSchoolData) (s *models.School, err *ce.Error)
}

type schoolRepository struct {
	database databases.SchoolDatabase
}

func NewSchoolRepository(db databases.SchoolDatabase) SchoolRepository {
	return &schoolRepository{database: db}
}

func (r *schoolRepository) Create(ctx context.Context, data *models.CreateSchoolData) (*models.School, *ce.Error) {
	return r.database.Create(ctx, data)
}
