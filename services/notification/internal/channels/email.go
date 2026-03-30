package channels

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/ritchieridanko/klasshub/services/notification/configs"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/logger"
	"github.com/ritchieridanko/klasshub/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/klasshub/services/notification/internal/models"
	"github.com/ritchieridanko/klasshub/services/notification/internal/templates"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils"
	"github.com/ritchieridanko/klasshub/services/notification/internal/utils/ce"
	"gopkg.in/gomail.v2"
)

type EmailChannel interface {
	SendWelcome(ctx context.Context, msg *models.WelcomeEmailMsg) (err *ce.Error)
}

type emailChannel struct {
	config   *configs.Client
	sender   string
	logoURL  string
	mailer   *mailer.Mailer
	template *template.Template
}

func NewEmailChannel(cfg *configs.Client, sender, logoURL string, m *mailer.Mailer) (EmailChannel, error) {
	tmpl, err := template.ParseFS(templates.Email, "*.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email channel: %w", err)
	}
	return &emailChannel{
		config:   cfg,
		sender:   sender,
		logoURL:  logoURL,
		mailer:   m,
		template: tmpl,
	}, nil
}

func (c *emailChannel) SendWelcome(ctx context.Context, msg *models.WelcomeEmailMsg) *ce.Error {
	// URL Generation
	url, err := utils.GenerateTokenizedURL(c.config.URL.Admin, "/auth/verify-email", msg.VerificationToken)
	if err != nil {
		return ce.NewError(ce.CodeURLGenerationFailed, err)
	}

	// Template Building
	body, buildErr := c.buildTemplate(
		"welcome",
		map[string]any{
			"Subject":   "Welcome Aboard!",
			"Recipient": msg.Recipient,
			"Title":     "Welcome to KlassHub",
			"URL":       url,
			"LogoURL":   c.logoURL,
			"Year":      time.Now().UTC().Year(),
		},
	)
	if buildErr != nil {
		return buildErr
	}

	// Message Composition
	m := c.composeMessage([]string{msg.Recipient}, "Welcome Aboard!", body.String())

	// Email Delivery
	return c.send(m)
}

func (c *emailChannel) buildTemplate(template string, data map[string]any) (bytes.Buffer, *ce.Error) {
	var buf bytes.Buffer
	if err := c.template.ExecuteTemplate(&buf, template, data); err != nil {
		return bytes.Buffer{}, ce.NewError(
			ce.CodeEmailTemplatingFailed,
			err,
			logger.NewField("email_template", template),
		)
	}
	return buf, nil
}

func (c *emailChannel) composeMessage(recipients []string, subject, body string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.sender)
	msg.SetHeader("To", recipients...)
	msg.SetHeader("Subject", utils.ToMIMEBase64(subject))
	msg.SetBody("text/plain", "Please view this email in an HTML-compatible client!")
	msg.AddAlternative("text/html", body)
	return msg
}

func (c *emailChannel) send(msg *gomail.Message) *ce.Error {
	if err := c.mailer.Send(msg); err != nil {
		return ce.NewError(ce.CodeEmailDeliveryFailed, err)
	}
	return nil
}
