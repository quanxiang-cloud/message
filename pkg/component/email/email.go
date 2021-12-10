package email

import (
	"context"
	"flag"

	"github.com/quanxiang-cloud/message/pkg/component/event"
	"gopkg.in/gomail.v2"
)

var (
	// Host represents the host of the SMTP server.
	host string
	// Port represents the port of the SMTP server.
	port int
	// Username is the username to use to authenticate to the SMTP server.
	username string
	// Password is the password to use to authenticate to the SMTP server.
	password string
	// Alias sender alias name.
	alias string
	// sender sender email.
	sender string
)

func Prepare() {
	flag.StringVar(&host, "email-host", "", " the host of the SMTP server")
	flag.IntVar(&port, "email-port", 0, "represents the port of the SMTP server")
	flag.StringVar(&username, "email-username", "", "the username to use to authenticate to the SMTP server")
	flag.StringVar(&password, "email-password", "", "the password to use to authenticate to the SMTP server")
	flag.StringVar(&alias, "email-alias", "", "sender alias name")
	flag.StringVar(&sender, "email-sender", "", "ender email")
}

func New(ctx context.Context) (*Email, error) {
	return &Email{
		dialer: gomail.NewDialer(host, port, username, password),
	}, nil
}

type Email struct {
	dialer *gomail.Dialer
}

func (e *Email) Scaffold(ctx context.Context, data event.Data) error {
	if data.EmailSpec == nil {
		return event.ErrDataIsNil
	}

	return e.Send(ctx, data.EmailSpec)
}

func (e *Email) Send(ctx context.Context, data *event.EmailSpec) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", sender, alias)
	m.SetHeader("To", data.To...)
	m.SetHeader("Subject", data.Title)
	if data.Content == "" {
		data.Content = "text/html"
	}
	m.SetBody(data.ContentType, data.Content)

	if err := e.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
