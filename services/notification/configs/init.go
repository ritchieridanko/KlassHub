package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ritchieridanko/klasshub/services/notification/internal/constants"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Client   Client   `mapstructure:"client"`
	Database Database `mapstructure:"database"`
	Broker   Broker   `mapstructure:"broker"`
	Tracer   Tracer   `mapstructure:"tracer"`
	Mailer   Mailer   `mapstructure:"mailer"`
}

type App struct {
	Name    string `mapstructure:"name"`
	LogoURL string `mapstructure:"logo_url"`
	Env     string
}

type Client struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	URL struct {
		Admin string
		LMS   string
	}
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

type Broker struct {
	Brokers string `mapstructure:"brokers"`

	Subscriber struct {
		AuthCreated struct {
			Name           string        `mapstructure:"name"`
			MaxBytes       int           `mapstructure:"max_bytes"`
			MaxWait        time.Duration `mapstructure:"max_wait"`
			CommitInterval time.Duration `mapstructure:"commit_interval"`
			ProcessTimeout time.Duration `mapstructure:"process_timeout"`
		} `mapstructure:"auth_created"`
	} `mapstructure:"subscriber"`
}

type Tracer struct {
	Addr string
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Mailer struct {
	Host    string        `mapstructure:"host"`
	Port    int           `mapstructure:"port"`
	User    string        `mapstructure:"user"`
	Pass    string        `mapstructure:"pass"`
	From    string        `mapstructure:"from"`
	Timeout time.Duration `mapstructure:"timeout"`
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
	cfg.Client.URL.Admin = fmt.Sprintf("http://admin.%s:%d", cfg.Client.Host, cfg.Client.Port)
	cfg.Client.URL.LMS = fmt.Sprintf("http://lms.%s:%d", cfg.Client.Host, cfg.Client.Port)
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

	constants.EventTopicAC = cfg.Broker.Subscriber.AuthCreated.Name

	return &cfg, nil
}
