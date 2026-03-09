package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ritchieridanko/klasshub/services/auth/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Auth     Auth     `mapstructure:"auth"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	Cache    Cache    `mapstructure:"cache"`
	Tracer   Tracer   `mapstructure:"tracer"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string
}

type Auth struct {
	BCrypt struct {
		Cost int `mapstructure:"cost"`
	} `mapstructure:"bcrypt"`

	JWT struct {
		Issuer   string        `mapstructure:"issuer"`
		Secret   string        `mapstructure:"secret"`
		Duration time.Duration `mapstructure:"duration"`
	} `mapstructure:"jwt"`

	Duration struct {
		Session time.Duration `mapstructure:"session"`
	} `mapstructure:"duration"`
}

type Server struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	Timeout struct {
		Shutdown time.Duration `mapstructure:"shutdown"`
	} `mapstructure:"timeout"`
}

type Database struct {
	DSN             string
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Pass            string        `mapstructure:"pass"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConns        int           `mapstructure:"max_conns"`
	MinConns        int           `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

type Cache struct {
	Addr            string
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Pass            string        `mapstructure:"pass"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxActiveConns  int           `mapstructure:"max_active_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`

	Timeout struct {
		Dial  time.Duration `mapstructure:"dial"`
		Read  time.Duration `mapstructure:"read"`
		Write time.Duration `mapstructure:"write"`
	} `mapstructure:"timeout"`
}

type Tracer struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Init(path string) (*Config, error) {
	env := utils.NormalizeString(os.Getenv("APP_ENV"))
	if env == "" || (env != "dev" && env != "prod") {
		env = "dev"
	}
	if path == "" {
		path = "./configs"
	}

	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	cfg.App.Env = env
	cfg.Server.Addr = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	cfg.Cache.Addr = fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port)
	cfg.Tracer.Addr = fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)
	cfg.Database.DSN = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	return &cfg, nil
}
