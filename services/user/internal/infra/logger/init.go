package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Init(env string) (*zap.Logger, error) {
	var cfg zap.Config
	if env == "prod" {
		cfg = zap.NewProductionConfig()
		cfg.DisableStacktrace = true
		cfg.DisableCaller = false
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.DisableStacktrace = false
		cfg.DisableCaller = false
	}

	l, err := cfg.Build(zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	l.Sugar().Infof("[LOGGER] initialized (env=%s, level=%s)", env, l.Level().String())
	return l, nil
}
