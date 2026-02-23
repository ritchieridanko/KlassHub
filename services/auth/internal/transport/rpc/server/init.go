package server

import (
	"context"
	"fmt"
	"net"

	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"google.golang.org/grpc"
)

type Server struct {
	name   string
	config *configs.Server
	server *grpc.Server
	logger *logger.Logger
	ah     *handlers.AuthHandler
}

func Init(name string, cfg *configs.Server, l *logger.Logger, ah *handlers.AuthHandler) *Server {
	srv := grpc.NewServer()
	apis.RegisterAuthServiceServer(srv, ah)

	return &Server{
		name:   name,
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
