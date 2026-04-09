package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ritchieridanko/klasshub/services/account/configs"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/account/internal/transport/rpc/interceptors"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/validator"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Server struct {
	config *configs.Server
	server *grpc.Server
	logger *logger.Logger
	ah     *handlers.AccountHandler
}

func Init(cfg *configs.Server, name string, v *validator.Validator, l *logger.Logger, ah *handlers.AccountHandler) *Server {
	srv := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
		grpc.ChainUnaryInterceptor(
			interceptors.Request(),
			interceptors.Recovery(l),
			interceptors.Logging(l),
			interceptors.Auth(v),
		),
	)

	apis.RegisterAccountServiceServer(srv, ah)

	return &Server{
		config: cfg,
		server: srv,
		logger: l,
		ah:     ah,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("failed to build listener: %w", err)
	}
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.logger.Log("[SERVER] is running (host=%s, port=%d)", s.config.Host, s.config.Port)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return fmt.Errorf("failed to shutdown server: %w", ctx.Err())
	case <-stopped:
		return nil
	}
}
