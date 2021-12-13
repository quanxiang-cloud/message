package lowcode

import (
	"context"
	"net/http"

	"git.internal.yunify.com/qxp/misc/client"
)

type Lowcode struct {
	warden Warden
}

func NewLowcodeAuth() *Lowcode {
	return &Lowcode{
		warden: NewWarden(client.Config{
			Timeout:      20,
			MaxIdleConns: 10,
		}),
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
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	ctx := context.Background()

	profile, err := l.warden.CheckToken(ctx, token, authURL)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	r.Header.Add("Id", profile.UserID)
	return true
}
