package handlers

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/school/internal/models"
	"github.com/ritchieridanko/klasshub/services/school/internal/usecases"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type SchoolHandler struct {
	apis.UnimplementedSchoolServiceServer
	su usecases.SchoolUsecase
}

func NewSchoolHandler(su usecases.SchoolUsecase) *SchoolHandler {
	return &SchoolHandler{su: su}
}

func (h *SchoolHandler) GetSchoolID(ctx context.Context, req *apis.GetSchoolIDRequest) (*apis.GetSchoolIDResponse, error) {
	schoolID, err := h.su.GetSchoolID(
		ctx,
		&models.GetSchoolIDRequest{
			AuthID: req.GetAuthId(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &apis.GetSchoolIDResponse{SchoolId: schoolID}, nil
}
