package di

import (
	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/cache"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/klasshub/services/auth/internal/transport/rpc/handlers"
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

	acp *publisher.Publisher

	ar repositories.AuthRepository
	sr repositories.SessionRepository
	tr repositories.TokenRepository

	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator

	su usecases.SessionUsecase
	au usecases.AuthUsecase

	ah     *handlers.AuthHandler
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
	au := usecases.NewAuthUsecase(cfg.App.Name, cfg.Auth.Duration.Verification, su, ar, tr, tx, acp, v, b, l)

	// Handlers
	ah := handlers.NewAuthHandler(au)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, l, ah)

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
		ar:         ar,
		sr:         sr,
		tr:         tr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		su:         su,
		au:         au,
		ah:         ah,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
