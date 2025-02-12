package auth

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Authenticator interface {
	Authorize(username, password string) (string, error)
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,alphanum"`
}

type AuthResponseOK struct {
	Token string `json:"token"`
}

type AuthResponseError struct {
	Error string `json:"errors"`
}

func New(log *slog.Logger, authenticator Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req AuthRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("error decoding request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, AuthResponseError{
				Error: "error decoding request body: " + err.Error(),
			})

			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", validateErr.Error()))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, AuthResponseError{
				Error: validateErr.Error(),
			})

			return
		}

		token, err := authenticator.Authorize(req.Username, req.Password)
		if err != nil {

			log.Error("error authenticating user", slog.String("err", err.Error()))

			render.Status(r, http.StatusUnauthorized)

			render.JSON(w, r, AuthResponseError{
				Error: err.Error(),
			})

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, AuthResponseOK{
			Token: token,
		})
	}
}
