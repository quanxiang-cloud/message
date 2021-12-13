package lowcode

import (
	"context"
	"errors"
	"net/http"

	"git.internal.yunify.com/qxp/misc/client"
	"git.internal.yunify.com/qxp/misc/logger"
)

// Profile profile
type Profile struct {
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
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Access-Token", token)

	// header 封装requestID
	logger.HeadAdd(&req.Header, logger.STDRequestID(ctx).String)

	response, err := o.client.Do(req)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	return &Profile{
		UserID:   response.Header.Get("User-Id"),
		UserName: response.Header.Get("User-Name"),
	}, nil
}
