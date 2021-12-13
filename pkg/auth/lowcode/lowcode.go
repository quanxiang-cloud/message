package lowcode

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/pkg/client"
)

type Lowcode struct {
	warden Warden
	log    logr.Logger
}

func NewLowcodeAuth(log logr.Logger) *Lowcode {

	return &Lowcode{
		warden: NewWarden(client.Config{
			Timeout:      20,
			MaxIdleConns: 10,
		}),
		log: log.WithName("lowcode"),
	}
}

const authURL = "http://jwt/api/v1/jwt/check"

func (l *Lowcode) Auth(w http.ResponseWriter, r *http.Request) bool {
	token := r.URL.Query().Get("token")
	if token == "" {
		// try get token from header
		token = r.Header.Get("Access-Token")
	}

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		l.log.Info("can not get token")
		return false
	}

	ctx := context.Background()

	profile, err := l.warden.CheckToken(ctx, token, authURL)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		l.log.Error(err, "checkToken")
		return false
	}

	r.Header.Add("Id", profile.UserID)
	return true
}
