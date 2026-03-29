package mailer

import (
	"sync"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	mutex  sync.Mutex
	dialer *gomail.Dialer
	sender gomail.SendCloser
}

func NewMailer(dlr *gomail.Dialer) *Mailer {
	return &Mailer{dialer: dlr}
}

func (m *Mailer) Send(msg *gomail.Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.ensureConnection(); err != nil {
		return err
	}
	return gomail.Send(m.sender, msg)
}

func (m *Mailer) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.sender != nil {
		return m.sender.Close()
	}
	return nil
}

func (m *Mailer) ensureConnection() error {
	if m.sender == nil {
		sc, err := m.dialer.Dial()
		if err != nil {
			return err
		}

		m.sender = sc
	}
	return nil
}
