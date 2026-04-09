package services

import (
	"fmt"
	"strings"

	"github.com/ritchieridanko/klasshub/services/account/configs"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SchoolService struct {
	conn   *grpc.ClientConn
	client apis.SchoolServiceClient
}

func NewSchoolService(cfg *configs.Service, l *zap.Logger) (*SchoolService, error) {
	conn, err := grpc.NewClient(
		cfg.School.Addr,
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to school service: %w", err)
	}

	l.Sugar().Infof(
		"[%s] connected (host=%s, port=%d)",
		strings.ToUpper(cfg.School.Name), cfg.School.Host, cfg.School.Port,
	)
	return &SchoolService{
		conn:   conn,
		client: apis.NewSchoolServiceClient(conn),
	}, nil
}

func (s *SchoolService) Client() apis.SchoolServiceClient {
	return s.client
}

func (s *SchoolService) Close() error {
	return s.conn.Close()
}
