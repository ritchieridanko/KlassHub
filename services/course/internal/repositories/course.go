package repositories

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/course/internal/models"
	"github.com/ritchieridanko/klasshub/services/course/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils/ce"
)

type CourseRepository interface {
	Create(ctx context.Context, data *models.CreateCourseData) (c *models.Course, err *ce.Error)
}

type courseRepository struct {
	database databases.CourseDatabase
}

func NewCourseRepository(db databases.CourseDatabase) CourseRepository {
	return &courseRepository{database: db}
}

func (r *courseRepository) Create(ctx context.Context, data *models.CreateCourseData) (*models.Course, *ce.Error) {
	return r.database.Create(ctx, data)
}
