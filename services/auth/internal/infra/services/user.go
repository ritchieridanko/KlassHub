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

type UserServiceClient struct {
	conn *grpc.ClientConn
	usc  apis.UserServiceClient
}

func NewUserServiceClient(cfg *configs.Service, l *zap.Logger) (*UserServiceClient, error) {
	conn, err := grpc.NewClient(
		cfg.User.Addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user service client: %w", err)
	}

	l.Sugar().Infof(
		"[%s] is running (host=%s, port=%d)",
		strings.ToUpper(cfg.User.Name), cfg.User.Host, cfg.User.Port,
	)
	return &UserServiceClient{
		conn: conn,
		usc:  apis.NewUserServiceClient(conn),
	}, nil
}

func (s *UserServiceClient) Client() apis.UserServiceClient {
	return s.usc
}

func (s *UserServiceClient) Close() error {
	return s.conn.Close()
}
