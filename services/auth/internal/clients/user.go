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

type UserService interface {
	GetSchoolAndRole(ctx context.Context, authID int64) (schoolID int64, role string, err *ce.Error)
}

type userService struct {
	client apis.UserServiceClient
}

func NewUserService(usc apis.UserServiceClient) UserService {
	return &userService{client: usc}
}

func (s *userService) GetSchoolAndRole(ctx context.Context, authID int64) (int64, string, *ce.Error) {
	resp, err := s.client.GetSchoolAndRole(
		ctx,
		&apis.GetSchoolAndRoleRequest{AuthId: authID},
	)
	if err != nil {
		authIDField := logger.NewField("auth_id", authID)

		st, ok := status.FromError(err)
		if !ok {
			return 0, "", ce.NewError(
				ce.CodeUnknown,
				ce.MsgInternalServer,
				fmt.Errorf("user service: %w", err),
				authIDField,
			)
		}
		if st.Code() == codes.NotFound {
			return 0, "", ce.NewError(
				ce.CodeInvariantViolation,
				ce.MsgInternalServer,
				errors.New("user service: auth exists, but user missing"),
				authIDField,
			)
		}
		return 0, "", ce.NewError(
			ce.CodeInternal,
			st.Message(),
			fmt.Errorf("user service: %s", st.Message()),
			authIDField,
		)
	}
	return resp.GetSchoolId(), resp.GetRole(), nil
}
