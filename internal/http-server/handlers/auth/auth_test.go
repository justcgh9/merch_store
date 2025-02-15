package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/justcgh9/merch_store/internal/http-server/handlers/auth"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	mockAuth := mocks.NewAuthenticator(t)
	logger := slog.Default()
	handler := auth.New(logger, mockAuth)

	t.Run("successful authentication", func(t *testing.T) {
		mockAuth.On("Authorize", "validUser", "validPass").Return("validToken", nil).Once()

		reqBody, _ := json.Marshal(auth.AuthRequest{
			Username: "validUser",
			Password: "validPass",
		})
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer([]byte(`{"invalidJson":}`)))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody, _ := json.Marshal(auth.AuthRequest{}) // Missing fields
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("authentication error", func(t *testing.T) {
		mockAuth.On("Authorize", "invalidUser", "invalidPass").Return("", errors.New("authentication failed")).Once()

		reqBody, _ := json.Marshal(auth.AuthRequest{
			Username: "invalidUser",
			Password: "invalidPass",
		})
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
