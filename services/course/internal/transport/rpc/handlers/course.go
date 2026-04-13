package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/course/internal/models"
	"github.com/ritchieridanko/klasshub/services/course/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type CourseHandler struct {
	apis.UnimplementedCourseServiceServer
	cu usecases.CourseUsecase
}

func NewCourseHandler(cu usecases.CourseUsecase) *CourseHandler {
	return &CourseHandler{cu: cu}
}

func (h *CourseHandler) CreateCourse(ctx context.Context, req *apis.CreateCourseRequest) (*apis.CreateCourseResponse, error) {
	c, err := h.cu.CreateCourse(
		ctx,
		&models.CreateCourseReq{
			SchoolCourseID: req.SchoolCourseId,
			Name:           req.GetName(),
			Description:    req.Description,
			CoursePicture:  req.CoursePicture,
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.CreateCourseResponse{
		Course: h.toCourse(c),
	}, nil
}

func (h *CourseHandler) toCourse(c *models.Course) *apis.Course {
	if c == nil {
		return nil
	}
	return &apis.Course{
		Id:             c.ID.String(),
		SchoolCourseId: c.SchoolCourseID,
		Name:           c.Name,
		Description:    c.Description,
		CoursePicture:  c.CoursePicture,
		CreatedAt:      utils.ToTimestamp(&c.CreatedAt),
		UpdatedAt:      utils.ToTimestamp(&c.UpdatedAt),
	}
}
