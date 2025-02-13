package email

import (
	"github.com/devbydaniel/release-notes-go/config"
	"gopkg.in/gomail.v2"
)

type Config struct {
	To      string
	Subject string
	Body    string
}

func Send(e *Config) error {
	conf := config.New()
	server := conf.Email.Server
	port := conf.Email.Port
	user := conf.Email.User
	password := conf.Email.Password

	m := gomail.NewMessage()
	m.SetHeader("From", conf.Email.DefaultFrom)
	m.SetHeader("To", e.To)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/html", e.Body)
	d := gomail.NewDialer(server, port, user, password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
