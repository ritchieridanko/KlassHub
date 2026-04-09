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
	config    *configs.Config
	logger    *logger.Logger
	auc       clients.AuthClient
	acc       clients.AccountClient
	cookie    *cookie.Cookie
	jwt       *jwt.JWT
	validator *validator.Validator
	auh       *handlers.AuthHandler
	ach       *handlers.AccountHandler
	router    *router.Router
	server    *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(inf.Logger())

	// Clients
	auc := clients.NewAuthClient(inf.AuthService())
	acc := clients.NewAccountClient(inf.AccountService())

	// Utils
	c := cookie.Init(cfg.App.Env, "")
	j := jwt.Init(cfg.JWT.Secret)
	v := validator.Init()

	// Handlers
	auh := handlers.NewAuthHandler(auc, v, c)
	ach := handlers.NewAccountHandler(acc, c)

	// Router
	r := router.Init(&cfg.Client, cfg.App.Name, j, l, auh, ach)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:    cfg,
		logger:    l,
		auc:       auc,
		acc:       acc,
		cookie:    c,
		jwt:       j,
		validator: v,
		auh:       auh,
		ach:       ach,
		router:    r,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
