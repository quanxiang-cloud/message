package lowcode

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/message/pkg/client"
)

// Profile profile
type Profile struct {
	Code     int
	UserID   string
	UserName string
}

// Oauth2s oauth2s
type Warden interface {
	CheckToken(ctx context.Context, token, checkURI string) (*Profile, error)
}

type warden struct {
	client http.Client
}

// NewWarden
func NewWarden(conf client.Config) Warden {
	return &warden{
		client: client.New(conf),
	}
}

func (o *warden) CheckToken(ctx context.Context, token, checkURI string) (*Profile, error) {
	req, err := http.NewRequest("POST", checkURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Access-Token", token)

	response, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return &Profile{
			Code: response.StatusCode,
		}, nil
	}

	return &Profile{
		UserID:   response.Header.Get("User-Id"),
		UserName: response.Header.Get("User-Name"),
	}, nil
}
