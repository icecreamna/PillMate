package mailer

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

type Mailer struct {
	host string
	port int
	user string
	pass string
	from string
	app  string
}

func New() (*Mailer, error) {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}
	m := &Mailer{
		host: os.Getenv("SMTP_HOST"),
		port: port,
		user: os.Getenv("SMTP_USERNAME"),
		pass: os.Getenv("SMTP_PASSWORD"),
		from: os.Getenv("SMTP_FROM"),
		app:  os.Getenv("APP_NAME"),
	}
	if m.host == "" || m.user == "" || m.pass == "" || m.from == "" {
		return nil, fmt.Errorf("missing smtp envs")
	}
	return m, nil
}

func (m *Mailer) Send(to, subject, htmlBody, textBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	if htmlBody != "" {
		msg.SetBody("text/html", htmlBody)
		if textBody != "" {
			msg.AddAlternative("text/plain", textBody)
		}
	} else {
		msg.SetBody("text/plain", textBody)
	}
	d := gomail.NewDialer(m.host, m.port, m.user, m.pass)
	return d.DialAndSend(msg)
}
