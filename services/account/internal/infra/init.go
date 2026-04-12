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
	us     *services.UserService
	acp    *kafka.Writer
	asufp  *kafka.Writer
	ucfp   *kafka.Writer
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
	us, err := services.NewUserService(&cfg.Service, l)
	if err != nil {
		return nil, err
	}

	// Publishers
	acp := publisher.Init(
		cfg.Broker.Brokers,
		cfg.Broker.Publisher.AC.Name,
		&kafka.Murmur2Balancer{
			Consistent: true,
		},
		cfg.Broker.Publisher.AC.BatchSize,
		cfg.Broker.Publisher.AC.BatchTimeout,
		l,
	)
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
	ucfp := publisher.Init(
		cfg.Broker.Brokers,
		cfg.Broker.Publisher.UCF.Name,
		&kafka.Murmur2Balancer{
			Consistent: true,
		},
		cfg.Broker.Publisher.UCF.BatchSize,
		cfg.Broker.Publisher.UCF.BatchTimeout,
		l,
	)

	return &Infra{
		config: cfg,
		logger: l,
		tracer: t,
		as:     as,
		ss:     ss,
		us:     us,
		acp:    acp,
		asufp:  asufp,
		ucfp:   ucfp,
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

func (i *Infra) UserService() apis.UserServiceClient {
	return i.us.Client()
}

func (i *Infra) PublisherAC() *kafka.Writer {
	return i.acp
}

func (i *Infra) PublisherASUF() *kafka.Writer {
	return i.asufp
}

func (i *Infra) PublisherUCF() *kafka.Writer {
	return i.ucfp
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
	if err := i.us.Close(); err != nil {
		return fmt.Errorf("failed to close user service connection: %w", err)
	}
	if err := i.acp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", constants.EventTopicAC, err)
	}
	if err := i.asufp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", constants.EventTopicASUF, err)
	}
	if err := i.ucfp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", constants.EventTopicUCF, err)
	}
	return nil
}
