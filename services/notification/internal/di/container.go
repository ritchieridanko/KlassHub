package di

import (
	"fmt"

	"github.com/ritchieridanko/klasshub/services/notification/configs"
	"github.com/ritchieridanko/klasshub/services/notification/internal/channels"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/notification/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/notification/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/klasshub/services/notification/internal/transport/event/middlewares"
	"github.com/ritchieridanko/klasshub/services/notification/internal/usecases"
)

type Container struct {
	config   *configs.Config
	database *database.Database
	logger   *logger.Logger
	mailer   *mailer.Mailer

	acs *subscriber.Subscriber

	ec channels.EmailChannel

	edb databases.EventDatabase

	er repositories.EventRepository

	au usecases.AuthUsecase

	rqm handlers.Middleware
	rvm handlers.Middleware
	tm  handlers.Middleware
	lm  handlers.Middleware

	ah  *handlers.AuthHandler
	ach handlers.Handler
}

func Init(cfg *configs.Config, inf *infra.Infra) (*Container, error) {
	// Infra
	db := database.NewDatabase(inf.Database())
	l := logger.NewLogger(inf.Logger())
	m := mailer.NewMailer(inf.Mailer())

	// Subscribers
	acs := subscriber.NewSubscriber(cfg.Broker.Subscriber.AuthCreated.ProcessTimeout, inf.SubscriberAC(), l)

	// Channels
	ec, err := channels.NewEmailChannel(&cfg.Client, cfg.Mailer.From, cfg.App.LogoURL, m)
	if err != nil {
		return nil, err
	}

	// Databases
	edb := databases.NewEventDatabase(db)

	// Repositories
	er := repositories.NewEventRepository(edb)

	// Usecases
	au := usecases.NewAuthUsecase(cfg.App.Name, cfg.Mailer.Timeout, ec, er, l)

	// Middlewares
	rqm := middlewares.Request()
	rvm := middlewares.Recovery(l)
	tm := middlewares.Tracing()
	lm := middlewares.Logging(l)

	// Handlers
	ah := handlers.NewAuthHandler(au)
	ach := handlers.NewHandler(ah.OnAuthCreated, rqm, rvm, tm, lm)

	return &Container{
		config:   cfg,
		database: db,
		logger:   l,
		mailer:   m,
		acs:      acs,
		ec:       ec,
		edb:      edb,
		er:       er,
		au:       au,
		rqm:      rqm,
		rvm:      rvm,
		tm:       tm,
		lm:       lm,
		ah:       ah,
		ach:      ach,
	}, nil
}

func (c *Container) SubscriberAC() *subscriber.Subscriber {
	return c.acs
}

func (c *Container) HandlerAC() handlers.Handler {
	return c.ach
}

func (c *Container) Close() error {
	if err := c.mailer.Close(); err != nil {
		return fmt.Errorf("failed to close mailer: %w", err)
	}
	return nil
}
