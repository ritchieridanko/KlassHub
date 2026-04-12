package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ritchieridanko/klasshub/services/account/internal/constants"
	"github.com/ritchieridanko/klasshub/services/account/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	App     App     `mapstructure:"app"`
	Server  Server  `mapstructure:"server"`
	Service Service `mapstructure:"service"`
	Broker  Broker  `mapstructure:"broker"`
	Tracer  Tracer  `mapstructure:"tracer"`
}

type App struct {
	Name string `mapstructure:"name"`
	Env  string
}

type Server struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	Timeout struct {
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

	School struct {
		Name string `mapstructure:"name"`
		Addr string
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"school"`

	User struct {
		Name string `mapstructure:"name"`
		Addr string
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"user"`
}

type Broker struct {
	Brokers string `mapstructure:"brokers"`

	Publisher struct {
		AC struct {
			Name         string        `mapstructure:"name"`
			BatchSize    int           `mapstructure:"batch_size"`
			BatchTimeout time.Duration `mapstructure:"batch_timeout"`
		} `mapstructure:"auth_created"`

		ASUF struct {
			Name         string        `mapstructure:"name"`
			BatchSize    int           `mapstructure:"batch_size"`
			BatchTimeout time.Duration `mapstructure:"batch_timeout"`
		} `mapstructure:"auth_school_update_failed"`

		UCF struct {
			Name         string        `mapstructure:"name"`
			BatchSize    int           `mapstructure:"batch_size"`
			BatchTimeout time.Duration `mapstructure:"batch_timeout"`
		} `mapstructure:"user_creation_failed"`
	} `mapstructure:"publisher"`
}

type Tracer struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Init(path string) (*Config, error) {
	env := utils.NormalizeString(os.Getenv("APP_ENV"))
	if env != "dev" && env != "prod" {
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
	cfg.Service.School.Addr = fmt.Sprintf("%s:%d", cfg.Service.School.Host, cfg.Service.School.Port)
	cfg.Service.User.Addr = fmt.Sprintf("%s:%d", cfg.Service.User.Host, cfg.Service.User.Port)
	cfg.Tracer.Addr = fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)

	constants.EventTopicAC = cfg.Broker.Publisher.AC.Name
	constants.EventTopicASUF = cfg.Broker.Publisher.ASUF.Name
	constants.EventTopicUCF = cfg.Broker.Publisher.UCF.Name

	return &cfg, nil
}
