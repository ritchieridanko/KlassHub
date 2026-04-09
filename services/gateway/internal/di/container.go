package di

import (
	"github.com/ritchieridanko/klasshub/services/gateway/configs"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/clients"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/handlers"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/router"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/transport/http/server"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/cookie"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/jwt"
	"github.com/ritchieridanko/klasshub/services/gateway/internal/utils/validator"
)

type Container struct {
	config *configs.Config
	logger *logger.Logger

	auc clients.AuthClient
	acc clients.AccountClient
	sc  clients.SchoolClient

	cookie    *cookie.Cookie
	jwt       *jwt.JWT
	validator *validator.Validator

	auh *handlers.AuthHandler
	ach *handlers.AccountHandler
	sh  *handlers.SchoolHandler

	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(inf.Logger())

	// Clients
	auc := clients.NewAuthClient(inf.AuthService())
	acc := clients.NewAccountClient(inf.AccountService())
	sc := clients.NewSchoolClient(inf.SchoolService())

	// Utils
	c := cookie.Init(cfg.App.Env, "")
	j := jwt.Init(cfg.JWT.Secret)
	v := validator.Init()

	// Handlers
	auh := handlers.NewAuthHandler(auc, v, c)
	ach := handlers.NewAccountHandler(acc, c)
	sh := handlers.NewSchoolHandler(sc)

	// Router
	r := router.Init(&cfg.Client, cfg.App.Name, j, l, auh, ach, sh)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:    cfg,
		logger:    l,
		auc:       auc,
		acc:       acc,
		sc:        sc,
		cookie:    c,
		jwt:       j,
		validator: v,
		auh:       auh,
		ach:       ach,
		sh:        sh,
		router:    r,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
