package send

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/justcgh9/merch_store/internal/models/user"
)

type Sender interface {
	Send(from, to string, amount int) error
}

type SendRequest struct {
	To     string `json:"toUser" validate:"required, alphanum"`
	Amount int    `json:"amount" validate:"required, number"`
}

type SendResponseError struct {
	Error string `json:"errors"`
}

func New(log *slog.Logger, sender Sender) http.HandlerFunc {
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
			render.JSON(w, r, SendResponseError{
				Error: "could not get user info",
			})
			return
		}

		log = log.With(
			slog.String("username", userDTO.Username),
		)

		var req SendRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("error decoding request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, SendResponseError{
				Error: "error decoding request body: " + err.Error(),
			})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("err", validateErr.Error()))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, SendResponseError{
				Error: validateErr.Error(),
			})

			return
		}

		if err := sender.Send(userDTO.Username, req.To, req.Amount); err != nil {

			log.Error("error sending money", slog.String("err", err.Error()))

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, SendResponseError{
				Error: err.Error(),
			})

			return
		}

		render.Status(r, http.StatusOK)
	}
}
