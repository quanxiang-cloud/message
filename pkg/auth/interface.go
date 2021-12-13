package auth

import "net/http"

type Interface interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
}
