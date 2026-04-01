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
	ac        clients.AuthClient
	cookie    *cookie.Cookie
	jwt       *jwt.JWT
	validator *validator.Validator
	ah        *handlers.AuthHandler
	router    *router.Router
	server    *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(inf.Logger())

	// Clients
	ac := clients.NewAuthClient(inf.AuthService())

	// Utils
	c := cookie.Init(cfg.App.Env, "")
	j := jwt.Init(cfg.JWT.Secret)
	v := validator.Init()

	// Handlers
	ah := handlers.NewAuthHandler(ac, v, c)

	// Router
	r := router.Init(&cfg.Client, cfg.App.Name, j, l, ah)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:    cfg,
		logger:    l,
		ac:        ac,
		cookie:    c,
		jwt:       j,
		validator: v,
		ah:        ah,
		router:    r,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
