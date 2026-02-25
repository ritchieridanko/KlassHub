package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/services"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/tracer"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	database *pgxpool.Pool
	logger   *zap.Logger
	tracer   *tracer.Tracer

	usc *services.UserServiceClient
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		return nil, err
	}

	t, err := tracer.Init(cfg.App.Name, cfg.Tracer.Endpoint, l)
	if err != nil {
		return nil, err
	}

	usc, err := services.NewUserServiceClient(&cfg.Service, l)
	if err != nil {
		return nil, err
	}

	return &Infra{
		config:   cfg,
		database: db,
		logger:   l,
		tracer:   t,
		usc:      usc,
	}, nil
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) UserServiceClient() apis.UserServiceClient {
	return i.usc.Client()
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.usc.Close(); err != nil {
		return fmt.Errorf("failed to close user service client: %w", err)
	}

	i.database.Close()
	i.tracer.Cleanup()
	return nil
}
