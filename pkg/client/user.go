package client

import (
	"time"

	"git.internal.yunify.com/qxp/organizations/pkg/client"
	cabinclient "github.com/quanxiang-cloud/cabin/tailormade/client"
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
