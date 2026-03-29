package mailer

import (
	"github.com/ritchieridanko/klasshub/services/notification/configs"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

func Init(cfg *configs.Mailer, l *zap.Logger) *gomail.Dialer {
	dlr := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass)

	l.Sugar().Infof("[MAILER] initialized (host=%s, port=%d, from=%s)", cfg.Host, cfg.Port, cfg.From)
	return dlr
}
