package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/klasshub/services/course/configs"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/services"
	"github.com/ritchieridanko/klasshub/services/course/internal/infra/tracer"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	database *pgxpool.Pool
	logger   *zap.Logger
	tracer   *tracer.Tracer
	ss       *services.SchoolService
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

	t, err := tracer.Init(cfg.App.Env, cfg.App.Name, cfg.Tracer.Addr, l)
	if err != nil {
		return nil, err
	}

	// Services
	ss, err := services.NewSchoolService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}

	return &Infra{
		config:   cfg,
		database: db,
		logger:   l,
		tracer:   t,
		ss:       ss,
	}, nil
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) SchoolService() apis.SchoolServiceClient {
	return i.ss.Client()
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.ss.Close(); err != nil {
		return fmt.Errorf("failed to close school service connection: %w", err)
	}

	i.database.Close()
	return nil
}
