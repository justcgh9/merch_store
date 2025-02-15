package send_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/send"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/send/mocks"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/stretchr/testify/assert"
)

func TestSendHandler(t *testing.T) {
	mockSender := mocks.NewSender(t)
	logger := slog.Default()
	handler := send.New(logger, mockSender)

	t.Run("successful money transfer", func(t *testing.T) {
		mockSender.On("Send", "testUser", "anotherUser", 100).Return(nil).Once()

		reqBody := send.SendRequest{
			To:     "anotherUser",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(body))
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("send to self", func(t *testing.T) {
		reqBody := send.SendRequest{
			To:     "testUser",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(body))
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp send.SendResponseError
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "cannot send money to yourself", errResp.Error)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader([]byte("{invalid json}")))
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := send.SendRequest{
			To:     "",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(body))
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("sender returns error", func(t *testing.T) {
		mockSender.On("Send", "testUser", "anotherUser", 100).Return(errors.New("transfer failed")).Once()

		reqBody := send.SendRequest{
			To:     "anotherUser",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(body))
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp send.SendResponseError
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "transfer failed", errResp.Error)
	})

	t.Run("missing user in context", func(t *testing.T) {
		reqBody := send.SendRequest{
			To:     "anotherUser",
			Amount: 100,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/send", bytes.NewReader(body))
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var errResp send.SendResponseError
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "could not get user info", errResp.Error)
	})
}
