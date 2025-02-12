package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/justcgh9/merch_store/internal/models/user"
)

type Authenticator interface {
	Authenticate(token string) (user.UserDTO, error)
}

type AuthenticationError struct {
	Error string `json:"errors"`
}

const (
	authHeader = "Authorization"
	userDTOKey = "userDTO"
)

func New(log *slog.Logger, authenticator Authenticator) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.auth.New"

			authHeader := r.Header.Get(authHeader)

			if authHeader == "" {
				log.Error("missing authorization header")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, AuthenticationError{
					Error: "missing authorization header",
				})

				return
			}

			if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
				log.Error("invalid authorization header")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, AuthenticationError{
					Error: "invalid authorization header",
				})
				return
			}

			userDTO, err := authenticator.Authenticate(strings.Split(authHeader, "Bearer ")[1])
			if err != nil {
				log.Error("invalid jwt token", slog.String("err", err.Error()))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, AuthenticationError{
					Error: "invalid jwt token",
				})
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), userDTOKey, userDTO))

			next.ServeHTTP(w, r)
		}
	}
}