package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/klasshub/services/school/configs"
	"github.com/ritchieridanko/klasshub/services/school/internal/constants"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/tracer"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	database *pgxpool.Pool
	logger   *zap.Logger
	tracer   *tracer.Tracer
	asufs    *kafka.Reader
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

	// Subscribers
	asufs := subscriber.Init(
		cfg.App.Name,
		cfg.Broker.Brokers,
		cfg.Broker.Subscriber.ASUF.Name,
		cfg.Broker.Subscriber.ASUF.MaxBytes,
		cfg.Broker.Subscriber.ASUF.MaxWait,
		cfg.Broker.Subscriber.ASUF.CommitInterval,
		l,
	)

	return &Infra{
		config:   cfg,
		database: db,
		logger:   l,
		tracer:   t,
		asufs:    asufs,
	}, nil
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) SubscriberASUF() *kafka.Reader {
	return i.asufs
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.asufs.Close(); err != nil {
		return fmt.Errorf("failed to close subscriber (topic: %s): %w", constants.EventTopicASUF, err)
	}

	i.database.Close()
	return nil
}
