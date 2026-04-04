package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/klasshub/services/notification/configs"
	"github.com/ritchieridanko/klasshub/services/notification/internal/constants"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/tracer"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type Infra struct {
	config   *configs.Config
	database *pgxpool.Pool
	logger   *zap.Logger
	mailer   *gomail.Dialer
	tracer   *tracer.Tracer
	acs      *kafka.Reader
	avrs     *kafka.Reader
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

	m := mailer.Init(&cfg.Mailer, l)

	t, err := tracer.Init(cfg.App.Env, cfg.App.Name, cfg.Tracer.Addr, l)
	if err != nil {
		return nil, err
	}

	// Subscribers
	acs := subscriber.Init(
		cfg.App.Name,
		cfg.Broker.Brokers,
		cfg.Broker.Subscriber.AC.Name,
		cfg.Broker.Subscriber.AC.MaxBytes,
		cfg.Broker.Subscriber.AC.MaxWait,
		cfg.Broker.Subscriber.AC.CommitInterval,
		l,
	)
	avrs := subscriber.Init(
		cfg.App.Name,
		cfg.Broker.Brokers,
		cfg.Broker.Subscriber.AVR.Name,
		cfg.Broker.Subscriber.AVR.MaxBytes,
		cfg.Broker.Subscriber.AVR.MaxWait,
		cfg.Broker.Subscriber.AVR.CommitInterval,
		l,
	)

	return &Infra{
		config:   cfg,
		database: db,
		logger:   l,
		mailer:   m,
		tracer:   t,
		acs:      acs,
		avrs:     avrs,
	}, nil
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) Mailer() *gomail.Dialer {
	return i.mailer
}

func (i *Infra) SubscriberAC() *kafka.Reader {
	return i.acs
}

func (i *Infra) SubscriberAVR() *kafka.Reader {
	return i.avrs
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.acs.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber (topic: %s): %w", constants.EventTopicAC, err)
	}
	if err := i.avrs.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber (topic: %s): %w", constants.EventTopicAVR, err)
	}

	i.database.Close()
	return nil
}
