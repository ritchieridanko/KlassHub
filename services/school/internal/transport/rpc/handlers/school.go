package handlers

import (
	"context"

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

func (h *SchoolHandler) GetID(ctx context.Context, req *apis.GetIDRequest) (*apis.GetIDResponse, error) {
	schoolID, err := h.su.GetID(ctx, req.GetAuthId())
	if err != nil {
		return nil, err
	}
	return &apis.GetIDResponse{SchoolId: schoolID}, nil
}
