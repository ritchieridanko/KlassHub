package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App     App     `mapstructure:"app"`
	Client  Client  `mapstructure:"client"`
	Server  Server  `mapstructure:"server"`
	Service Service `mapstructure:"service"`
	Tracer  Tracer  `mapstructure:"tracer"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string
}

type Client struct {
	Hostname string `mapstructure:"hostname"`
	TLD      string `mapstructure:"tld"`
}

type Server struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	Timeout struct {
		Read     time.Duration `mapstructure:"read"`
		Write    time.Duration `mapstructure:"write"`
		Shutdown time.Duration `mapstructure:"shutdown"`
	} `mapstructure:"timeout"`
}

type Service struct {
	Auth struct {
		Name string `mapstructure:"name"`
		Addr string
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"auth"`
}

type Tracer struct {
	Endpoint string
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

func Init(path string) (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
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
	cfg.Service.Auth.Addr = fmt.Sprintf("%s:%d", cfg.Service.Auth.Host, cfg.Service.Auth.Port)
	cfg.Tracer.Endpoint = fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)

	return &cfg, nil
}
