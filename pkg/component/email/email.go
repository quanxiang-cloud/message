package email

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"io"
	"net/mail"
	"os"

	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/pkg/cache"
	"github.com/quanxiang-cloud/message/pkg/client"
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
	// CA cert file path.
	caCertPath string
)

const (
	//Certificate CA Cert Type
	Certificate = "CERTIFICATE"
)

// Prepare Prepare
func Prepare() {
	cache.PrepareCache()
	flag.StringVar(&host, "email-host", "", " the host of the SMTP server")
	flag.IntVar(&port, "email-port", 0, "represents the port of the SMTP server")
	flag.StringVar(&username, "email-username", "", "the username to use to authenticate to the SMTP server")
	flag.StringVar(&password, "email-password", "", "the password to use to authenticate to the SMTP server")
	flag.StringVar(&alias, "email-alias", "", "sender alias name")
	flag.StringVar(&sender, "email-sender", "", "ender email")
	flag.StringVar(&caCertPath, "ca-cert", "", "CA cert file path")
}

// New New
func New(ctx context.Context, log logr.Logger) (*Email, error) {
	attachCache, err := cache.NewCache(log)
	if err != nil {
		return nil, err
	}
	dialer, err := getMailDialer()
	if err != nil {
		return nil, err
	}
	return &Email{
		dialer:      dialer,
		log:         log.WithName("email"),
		attachCache: attachCache,
		fileServer:  client.NewFileServer(),
	}, nil
}

// Email Email
type Email struct {
	log         logr.Logger
	dialer      *gomail.Dialer
	attachCache cache.Cache
	fileServer  client.FileServerAPI
}

// Scaffold Scaffold
func (e *Email) Scaffold(ctx context.Context, data event.Data) error {
	if data.EmailSpec == nil {
		return event.ErrDataIsNil
	}
	err := e.warningData(ctx, data.EmailSpec)
	if err != nil {
		return err
	}
	return e.Send(ctx, data.EmailSpec)
}

func (e *Email) warningData(ctx context.Context, data *event.EmailSpec) error {
	toList := make([]string, 0, len(data.To))
	for _, t := range data.To {
		_, err := mail.ParseAddress(t)
		if err != nil {
			e.log.Error(err, "email address format error")
			continue
		}
		toList = append(toList, t)
	}
	data.To = toList
	return nil
}

// Send Send
func (e *Email) Send(ctx context.Context, data *event.EmailSpec) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", sender, alias)
	m.SetHeader("To", data.To...)
	m.SetHeader("Subject", data.Title)
	if data.ContentType == "" {
		data.ContentType = "text/html"
	}
	m.SetBody(data.ContentType, data.Content)
	for _, attach := range data.Attachments {
		content, err := e.getAttachment(ctx, attach.Path)
		if err != nil {
			e.log.Error(err, "Get Attach")
			continue
		}
		m.Attach(attach.Name, gomail.SetCopyFunc(func(writer io.Writer) error {
			_, err := writer.Write(content)
			return err
		}))
	}
	if err := e.dialer.DialAndSend(m); err != nil {
		e.log.Error(err, "DialAndSend")
		return err
	}
	return nil
}

func (e *Email) getAttachment(ctx context.Context, path string) ([]byte, error) {
	content, err := e.attachCache.Get(path)
	if err != nil {
		content, err = e.fileServer.GetFile(ctx, path)
		if err != nil {
			return nil, err
		}
		err = e.attachCache.Push(path, content)
		if err != nil {
			e.log.Error(err, "Cache invalidation")
		}
	}
	return content, nil
}

func getMailDialer() (*gomail.Dialer, error) {
	dialer := gomail.NewDialer(host, port, username, password)
	if caCertPath != "" {
		crtb, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, err
		}
		block, crtb := pem.Decode(crtb)
		if block == nil {
			return nil, errors.New("invalid PEM certificate")
		}
		if block.Type != Certificate {
			return nil, errors.New("cert not right")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		caPool := x509.NewCertPool()
		caPool.AddCert(cert)
		dialer.SSL = true
		dialer.TLSConfig.RootCAs = caPool
	}
	return dialer, nil
}
