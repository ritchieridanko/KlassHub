package infra

import (
	"fmt"

	"github.com/ritchieridanko/klasshub/services/gateway/configs"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/services"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/tracer"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"go.uber.org/zap"
)

type Infra struct {
	config *configs.Config
	logger *zap.Logger
	tracer *tracer.Tracer
	aus    *services.AuthService
	acs    *services.AccountService
	ss     *services.SchoolService
	us     *services.UserService
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	t, err := tracer.Init(cfg.App.Env, cfg.App.Name, cfg.Tracer.Addr, l)
	if err != nil {
		return nil, err
	}

	// Services
	aus, err := services.NewAuthService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}
	acs, err := services.NewAccountService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}
	ss, err := services.NewSchoolService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}
	us, err := services.NewUserService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}

	return &Infra{
		config: cfg,
		logger: l,
		tracer: t,
		aus:    aus,
		acs:    acs,
		ss:     ss,
		us:     us,
	}, nil
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) AuthService() apis.AuthServiceClient {
	return i.aus.Client()
}

func (i *Infra) AccountService() apis.AccountServiceClient {
	return i.acs.Client()
}

func (i *Infra) SchoolService() apis.SchoolServiceClient {
	return i.ss.Client()
}

func (i *Infra) UserService() apis.UserServiceClient {
	return i.us.Client()
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.aus.Close(); err != nil {
		return fmt.Errorf("failed to close auth service connection: %w", err)
	}
	if err := i.acs.Close(); err != nil {
		return fmt.Errorf("failed to close account service connection: %w", err)
	}
	if err := i.ss.Close(); err != nil {
		return fmt.Errorf("failed to close school service connection: %w", err)
	}
	if err := i.us.Close(); err != nil {
		return fmt.Errorf("failed to close user service connection: %w", err)
	}
	return nil
}
