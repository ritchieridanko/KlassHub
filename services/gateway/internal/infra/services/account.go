package services

import (
	"fmt"
	"strings"

	"github.com/ritchieridanko/klasshub/services/gateway/configs"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountService struct {
	conn   *grpc.ClientConn
	client apis.AccountServiceClient
}

func NewAccountService(cfg *configs.Service, l *zap.Logger) (*AccountService, error) {
	conn, err := grpc.NewClient(
		cfg.Account.Addr,
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to account service: %w", err)
	}

	l.Sugar().Infof(
		"[%s] connected (host=%s, port=%d)",
		strings.ToUpper(cfg.Account.Name), cfg.Account.Host, cfg.Account.Port,
	)
	return &AccountService{
		conn:   conn,
		client: apis.NewAccountServiceClient(conn),
	}, nil
}

func (s *AccountService) Client() apis.AccountServiceClient {
	return s.client
}

func (s *AccountService) Close() error {
	return s.conn.Close()
}
