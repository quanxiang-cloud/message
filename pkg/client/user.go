package client

import (
	"time"

	cabinclient "github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/organizations/pkg/client"
)

type User interface {
	client.User
}

func NewUser(conf Config) User {
	return client.NewUser(cabinclient.Config{
		Timeout:      time.Second * conf.Timeout,
		MaxIdleConns: conf.MaxIdleConns,
	})
}
