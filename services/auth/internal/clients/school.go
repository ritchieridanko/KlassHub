package clients

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SchoolService interface {
	GetID(ctx context.Context, authID int64) (schoolID int64, err *ce.Error)
}

type schoolService struct {
	client apis.SchoolServiceClient
}

func NewSchoolService(ssc apis.SchoolServiceClient) SchoolService {
	return &schoolService{client: ssc}
}

func (s *schoolService) GetID(ctx context.Context, authID int64) (int64, *ce.Error) {
	resp, err := s.client.GetID(
		ctx,
		&apis.GetIDRequest{AuthId: authID},
	)
	if err != nil {
		authIDField := logger.NewField("auth_id", authID)

		st, ok := status.FromError(err)
		if !ok {
			return 0, ce.NewError(
				ce.CodeUnknown,
				ce.MsgInternalServer,
				fmt.Errorf("school service: %w", err),
				authIDField,
			)
		}
		if st.Code() == codes.NotFound {
			return 0, ce.NewError(
				ce.CodeInvariantViolation,
				ce.MsgInternalServer,
				errors.New("school service: auth exists, but school missing"),
				authIDField,
			)
		}
		return 0, ce.NewError(
			ce.CodeInternal,
			st.Message(),
			fmt.Errorf("school service: %s", st.Message()),
			authIDField,
		)
	}
	return resp.GetSchoolId(), nil
}
