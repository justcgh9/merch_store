package buy

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/justcgh9/merch_store/internal/models/user"
)

type Buyer interface {
	Authorize(token string) (user.UserDTO, error)
	Buy(username, item string) error
}

type BuyResponseError struct {
	Error string `json:"errors"`
}

const (
	itemParam = "item"
	authHeader = "Authorization"
)

func New(log *slog.Logger, buyer Buyer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.buy.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		item := chi.URLParam(r, itemParam)

		authHeader := r.Header.Get(authHeader)

		if authHeader == "" {
			log.Error("missing authorization header")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, BuyResponseError{
				Error: "missing authorization header",
			})

			return
		}

		if len(authHeader) < 7 || authHeader[:6] != "Bearer" {
			log.Error("invalid authorization header")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, BuyResponseError{
				Error: "invalid authorization header",
			})
			return
		}

		userDTO, err := buyer.Authorize(strings.Split(authHeader, "Bearer ")[1])
		if err != nil {
			log.Error("invalid jwt token", slog.String("err", err.Error()))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, BuyResponseError{
				Error: "invalid jwt token",
			})
			return
		}

		log = log.With(
			slog.String("username", userDTO.Username),
		)

		err = buyer.Buy(userDTO.Username, item)
		if err != nil {
			// TODO distinguish the errors
			log.Error("could not buy " + item, slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, BuyResponseError{
				Error: "could not buy " + item,
			})
			return
		}

		render.Status(r, http.StatusOK)
	}
}