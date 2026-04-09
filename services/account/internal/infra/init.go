package infra

import (
	"fmt"

	"github.com/ritchieridanko/klasshub/services/account/configs"
	"github.com/ritchieridanko/klasshub/services/account/internal/constants"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/services"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/tracer"
	"github.com/ritchieridanko/klasshub/shared/contract/apis/v1"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Infra struct {
	config *configs.Config
	logger *zap.Logger
	tracer *tracer.Tracer
	as     *services.AuthService
	ss     *services.SchoolService
	asufp  *kafka.Writer
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
	as, err := services.NewAuthService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}
	ss, err := services.NewSchoolService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}

	// Publishers
	asufp := publisher.Init(
		cfg.Broker.Brokers,
		cfg.Broker.Publisher.ASUF.Name,
		&kafka.Murmur2Balancer{
			Consistent: true,
		},
		cfg.Broker.Publisher.ASUF.BatchSize,
		cfg.Broker.Publisher.ASUF.BatchTimeout,
		l,
	)

	return &Infra{
		config: cfg,
		logger: l,
		tracer: t,
		as:     as,
		ss:     ss,
		asufp:  asufp,
	}, nil
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) AuthService() apis.AuthServiceClient {
	return i.as.Client()
}

func (i *Infra) SchoolService() apis.SchoolServiceClient {
	return i.ss.Client()
}

func (i *Infra) PublisherASUF() *kafka.Writer {
	return i.asufp
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.as.Close(); err != nil {
		return fmt.Errorf("failed to close auth service connection: %w", err)
	}
	if err := i.ss.Close(); err != nil {
		return fmt.Errorf("failed to close school service connection: %w", err)
	}
	if err := i.asufp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", constants.EventTopicASUF, err)
	}
	return nil
}
