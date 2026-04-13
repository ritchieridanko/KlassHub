package clients

import (
	"context"

	"github.com/ritchieridanko/klasshub/services/course/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/course/internal/utils/ce"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
)

type SchoolClient interface {
	SchoolExists(ctx context.Context, schoolID int64) (exists bool, err *ce.Error)
}

type schoolClient struct {
	client apis.SchoolServiceClient
}

func NewSchoolClient(c apis.SchoolServiceClient) SchoolClient {
	return &schoolClient{client: c}
}

func (c *schoolClient) SchoolExists(ctx context.Context, schoolID int64) (bool, *ce.Error) {
	resp, err := c.client.SchoolExists(
		ctx,
		&apis.SchoolExistenceCheckRequest{
			SchoolId: schoolID,
		},
	)
	if err != nil {
		return false, ce.FromGRPCErr(
			err,
		).Append(
			logger.NewField("service", "school"),
		)
	}
	return resp.GetExists(), nil
}
