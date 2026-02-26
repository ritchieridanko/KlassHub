package services

import (
	"fmt"
	"strings"

	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SchoolServiceClient struct {
	conn *grpc.ClientConn
	ssc  apis.SchoolServiceClient
}

func NewSchoolServiceClient(cfg *configs.Service, l *zap.Logger) (*SchoolServiceClient, error) {
	conn, err := grpc.NewClient(
		cfg.School.Addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize school service client: %w", err)
	}

	l.Sugar().Infof(
		"[%s] is running (host=%s, port=%d)",
		strings.ToUpper(cfg.School.Name), cfg.School.Host, cfg.School.Port,
	)
	return &SchoolServiceClient{
		conn: conn,
		ssc:  apis.NewSchoolServiceClient(conn),
	}, nil
}

func (s *SchoolServiceClient) Client() apis.SchoolServiceClient {
	return s.ssc
}

func (s *SchoolServiceClient) Close() error {
	return s.conn.Close()
}
