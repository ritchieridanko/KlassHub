package databases

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/course/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/course/internal/models"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils/ce"
)

type CourseDatabase interface {
	Create(ctx context.Context, data *models.CreateCourseData) (c *models.Course, err *ce.Error)
}

type courseDatabase struct {
	database *database.Database
}

func NewCourseDatabase(db *database.Database) CourseDatabase {
	return &courseDatabase{database: db}
}

func (d *courseDatabase) Create(ctx context.Context, data *models.CreateCourseData) (*models.Course, *ce.Error) {
	query := `
		INSERT INTO courses (
			id, school_id, school_course_id,
			name, description, course_picture
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING
			id, school_course_id, name, description,
			course_picture, created_at, updated_at
	`

	var c models.Course
	err := d.database.Query(
		ctx, query,
		data.ID, data.SchoolID, data.SchoolCourseID,
		data.Name, data.Description, data.CoursePicture,
	).Scan(
		&c.ID, &c.SchoolCourseID, &c.Name, &c.Description,
		&c.CoursePicture, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create course: %w", err),
			logger.NewField("school_id", data.SchoolID),
		)
	}

	return &c, nil
}
