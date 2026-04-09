package di

import (
	"github.com/ritchieridanko/klasshub/services/school/configs"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/school/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/school/internal/repositories/databases"
	eHandlers "github.com/ritchieridanko/klasshub/services/school/internal/transport/event/handlers"
	eMiddlewares "github.com/ritchieridanko/klasshub/services/school/internal/transport/event/middlewares"
	rHandlers "github.com/ritchieridanko/klasshub/services/school/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/school/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/school/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/school/internal/utils/validator"
	"github.com/ritchieridanko/klasshub/shared/data"
)

type Container struct {
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	asufs *subscriber.Subscriber

	sd *data.School

	sdb databases.SchoolDatabase

	sr repositories.SchoolRepository

	validator *validator.Validator

	su usecases.SchoolUsecase

	// Event Middlewares
	rqm eHandlers.Middleware
	rvm eHandlers.Middleware
	tm  eHandlers.Middleware
	lm  eHandlers.Middleware

	// RPC Handlers
	sh *rHandlers.SchoolHandler

	// Event Handlers
	ah    *eHandlers.AuthHandler
	asufh eHandlers.Handler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra, sd *data.School) *Container {
	// Infra
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Subscribers
	asufs := subscriber.NewSubscriber(cfg.Broker.Subscriber.ASUF.ProcessTimeout, inf.SubscriberASUF(), l)

	// Databases
	sdb := databases.NewSchoolDatabase(db)

	// Repositories
	sr := repositories.NewSchoolRepository(sdb)

	// Utils
	v := validator.Init(sd)

	// Usecases
	su := usecases.NewSchoolUsecase(cfg.App.Name, sr, v)

	// Event Middlewares
	rqm := eMiddlewares.Request()
	rvm := eMiddlewares.Recovery(l)
	tm := eMiddlewares.Tracing()
	lm := eMiddlewares.Logging(l)

	// RPC Handlers
	sh := rHandlers.NewSchoolHandler(su)

	// Event Handlers
	ah := eHandlers.NewAuthHandler(su)
	asufh := eHandlers.NewHandler(ah.OnAuthSchoolUpdateFailed, rqm, rvm, tm, lm)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, v, l, sh)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		asufs:      asufs,
		sd:         sd,
		sdb:        sdb,
		sr:         sr,
		validator:  v,
		su:         su,
		rqm:        rqm,
		rvm:        rvm,
		tm:         tm,
		lm:         lm,
		sh:         sh,
		ah:         ah,
		asufh:      asufh,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}

func (c *Container) SubscriberASUF() *subscriber.Subscriber {
	return c.asufs
}

func (c *Container) HandlerASUF() eHandlers.Handler {
	return c.asufh
}
