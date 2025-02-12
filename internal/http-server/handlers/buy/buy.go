package buy

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/justcgh9/merch_store/internal/models/user"
)

type Buyer interface {
	Buy(username, item string) error
}

type BuyResponseError struct {
	Error string `json:"errors"`
}

const (
	itemParam  = "item"
	userDTOKey = "userDTO"
)

func New(log *slog.Logger, buyer Buyer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.buy.New"
		
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userDTO, ok := r.Context().Value(userDTOKey).(user.UserDTO)

		log = log.With(
			slog.String("username", userDTO.Username),
		)

		item := chi.URLParam(r, itemParam)

		if !ok {
			log.Error("could not get user info")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, BuyResponseError{
				Error: "could not get user info",
			})
			return
		}

		err := buyer.Buy(userDTO.Username, item)
		if err != nil {
			// TODO distinguish the errors
			log.Error("could not buy "+item, slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, BuyResponseError{
				Error: "could not buy " + item,
			})
			return
		}

		render.Status(r, http.StatusOK)
	}
}
