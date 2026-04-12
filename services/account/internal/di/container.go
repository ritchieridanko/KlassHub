package di

import (
	"github.com/ritchieridanko/klasshub/services/account/configs"
	"github.com/ritchieridanko/klasshub/services/account/internal/clients"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/account/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/account/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/account/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/account/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils/validator"
)

type Container struct {
	config *configs.Config
	logger *logger.Logger

	ac clients.AuthClient
	sc clients.SchoolClient
	uc clients.UserClient

	acp   *publisher.Publisher
	asufp *publisher.Publisher
	ucfp  *publisher.Publisher

	validator *validator.Validator

	au usecases.AccountUsecase

	ah     *handlers.AccountHandler
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(inf.Logger())

	// Clients
	ac := clients.NewAuthClient(inf.AuthService())
	sc := clients.NewSchoolClient(inf.SchoolService())
	uc := clients.NewUserClient(inf.UserService())

	// Publishers
	acp := publisher.NewPublisher(inf.PublisherAC())
	asufp := publisher.NewPublisher(inf.PublisherASUF())
	ucfp := publisher.NewPublisher(inf.PublisherUCF())

	// Utils
	v := validator.Init()

	// Usecases
	au := usecases.NewAccountUsecase(cfg.App.Name, ac, sc, uc, acp, asufp, ucfp, l)

	// Handlers
	ah := handlers.NewAccountHandler(au)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, v, l, ah)

	return &Container{
		config:    cfg,
		logger:    l,
		ac:        ac,
		sc:        sc,
		uc:        uc,
		acp:       acp,
		asufp:     asufp,
		ucfp:      ucfp,
		validator: v,
		au:        au,
		ah:        ah,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
