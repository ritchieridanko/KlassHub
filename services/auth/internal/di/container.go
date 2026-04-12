package di

import (
	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/subscriber"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/databases"
	eHandlers "github.com/ritchieridanko/klasshub/services/auth/internal/transport/event/handlers"
	eMiddlewares "github.com/ritchieridanko/klasshub/services/auth/internal/transport/event/middlewares"
	rHandlers "github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/server"
	"github.com/ritchieridanko/klasshub/services/auth/internal/usecases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/bcrypt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/jwt"
	"github.com/ritchieridanko/klasshub/services/auth/internal/utils/validator"
)

type Container struct {
	config     *configs.Config
	cache      *cache.Cache
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	adb databases.AuthDatabase
	sdb databases.SessionDatabase

	acc caches.AuthCache
	tcc caches.TokenCache

	acp  *publisher.Publisher
	avrp *publisher.Publisher

	ucfs *subscriber.Subscriber

	ar repositories.AuthRepository
	sr repositories.SessionRepository
	tr repositories.TokenRepository

	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator

	su usecases.SessionUsecase
	au usecases.AuthUsecase

	// Event Middlewares
	rqm eHandlers.Middleware
	rvm eHandlers.Middleware
	tm  eHandlers.Middleware
	lm  eHandlers.Middleware

	// RPC Handlers
	ah *rHandlers.AuthHandler

	// Event Handlers
	uh   *eHandlers.UserHandler
	ucfh eHandlers.Handler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	cc := cache.NewCache(inf.Cache())
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	adb := databases.NewAuthDatabase(db)
	sdb := databases.NewSessionDatabase(db)

	// Caches
	acc := caches.NewAuthCache(cc)
	tcc := caches.NewTokenCache(&cfg.Auth, cc)

	// Publishers
	acp := publisher.NewPublisher(inf.PublisherAC())
	avrp := publisher.NewPublisher(inf.PublisherAVR())

	// Subscribers
	ucfs := subscriber.NewSubscriber(cfg.Broker.Subscriber.UCF.ProcessTimeout, inf.SubscriberUCF(), l)

	// Repositories
	ar := repositories.NewAuthRepository(adb, acc)
	sr := repositories.NewSessionRepository(sdb)
	tr := repositories.NewTokenRepository(tcc)

	// Utils
	b := bcrypt.Init(cfg.Auth.BCrypt.Cost)
	j := jwt.Init(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := validator.Init()

	// Usecases
	su := usecases.NewSessionUsecase(cfg.App.Name, cfg.Auth.JWT.Duration, cfg.Auth.Duration.Session, sr, tx, v, j)
	au := usecases.NewAuthUsecase(cfg.App.Name, cfg.Auth.Duration.Verification, su, ar, tr, tx, acp, avrp, v, b, l)

	// Event Middlewares
	rqm := eMiddlewares.Request()
	rvm := eMiddlewares.Recovery(l)
	tm := eMiddlewares.Tracing()
	lm := eMiddlewares.Logging(l)

	// RPC Handlers
	ah := rHandlers.NewAuthHandler(au)

	// Event Handlers
	uh := eHandlers.NewUserHandler(au)
	ucfh := eHandlers.NewHandler(uh.OnUserCreationFailed, rqm, rvm, tm, lm)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, v, l, ah)

	return &Container{
		config:     cfg,
		cache:      cc,
		database:   db,
		transactor: tx,
		logger:     l,
		adb:        adb,
		sdb:        sdb,
		acc:        acc,
		tcc:        tcc,
		acp:        acp,
		avrp:       avrp,
		ucfs:       ucfs,
		ar:         ar,
		sr:         sr,
		tr:         tr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		su:         su,
		au:         au,
		rqm:        rqm,
		rvm:        rvm,
		tm:         tm,
		lm:         lm,
		ah:         ah,
		uh:         uh,
		ucfh:       ucfh,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}

func (c *Container) SubscriberUCF() *subscriber.Subscriber {
	return c.ucfs
}

func (c *Container) HandlerUCF() eHandlers.Handler {
	return c.ucfh
}
