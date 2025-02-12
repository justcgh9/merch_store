package info

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/justcgh9/merch_store/internal/models/inventory"
	"github.com/justcgh9/merch_store/internal/models/user"
)

type Informator interface {
	Informate(username string) (inventory.Info, error)
}

type InfoResponseError struct {
	Error string `json:"errors"`
}

type InfoResponseOk = inventory.Info

func New(log *slog.Logger, informator Informator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.send.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userDTO, ok := r.Context().Value(user.UserDTOKey).(user.UserDTO)

		if !ok {
			log.Error("could not get user info")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, InfoResponseError{
				Error: "could not get user info",
			})
			return
		}

		log = log.With(
			slog.String("username", userDTO.Username),
		)

		var resp InfoResponseOk
		resp, err := informator.Informate(userDTO.Username)
		if err != nil {
			log.Error("failed to get user info", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, InfoResponseError{
				Error: err.Error(),
			})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}
