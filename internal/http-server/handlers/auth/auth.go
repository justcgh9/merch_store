package auth

import (
	"log/slog"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,aplhanum"`
}

type Authenticator interface {

}

func New(log *slog.Logger, authenticator Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}