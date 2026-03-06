package di

import (
	"github.com/ritchieridanko/klasshub/services/auth/configs"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/database"
	"github.com/ritchieridanko/klasshub/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/auth/internal/repositories"
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
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	adb databases.AuthDatabase
	sdb databases.SessionDatabase

	ar repositories.AuthRepository
	sr repositories.SessionRepository

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
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	adb := databases.NewAuthDatabase(db)
	sdb := databases.NewSessionDatabase(db)

	// Repositories
	ar := repositories.NewAuthRepository(adb)
	sr := repositories.NewSessionRepository(sdb)

	// Utils
	b := bcrypt.Init(cfg.Auth.BCrypt.Cost)
	j := jwt.Init(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := validator.Init()

	// Usecases
	su := usecases.NewSessionUsecase(cfg.App.Name, cfg.Auth.JWT.Duration, cfg.Auth.Duration.Session, sr, tx, v, j)
	au := usecases.NewAuthUsecase(cfg.App.Name, su, ar, tx, v, b)

	// Handlers
	ah := handlers.NewAuthHandler(au)

	// Server
	srv := server.Init(&cfg.Server, cfg.App.Name, l, ah)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		adb:        adb,
		sdb:        sdb,
		ar:         ar,
		sr:         sr,
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
