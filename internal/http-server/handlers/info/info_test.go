package info_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/info"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/info/mocks"
	"github.com/justcgh9/merch_store/internal/models/inventory"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/stretchr/testify/assert"
)

func TestInfoHandler(t *testing.T) {
	mockInformator := mocks.NewInformator(t)
	logger := slog.Default()
	handler := info.New(logger, mockInformator)

	t.Run("successful information retrieval", func(t *testing.T) {
		expectedInfo := inventory.Info{
			Balance: 1000,
			Inventory: []inventory.Item{
				{Type: "t_shirt", Quantity: 1},
				{Type: "cup", Quantity: 2},
			},
		}

		mockInformator.On("Informate", "testUser").Return(expectedInfo, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got inventory.Info
		err := json.NewDecoder(resp.Body).Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, expectedInfo, got)
	})

	t.Run("informator error", func(t *testing.T) {
		mockInformator.On("Informate", "testUser").Return(inventory.Info{}, errors.New("failed")).Once()

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp info.InfoResponseError
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "failed", errResp.Error)
	})

	t.Run("missing user in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		chiCtx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var errResp info.InfoResponseError
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		assert.NoError(t, err)
		assert.Equal(t, "could not get user info", errResp.Error)
	})
}
