package buy_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/buy"
	"github.com/justcgh9/merch_store/internal/http-server/handlers/buy/mocks"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/stretchr/testify/assert"
)

func TestBuyHandler(t *testing.T) {
	logger := slog.Default()

	t.Run("successful purchase", func(t *testing.T) {		
		mockBuyer := mocks.NewBuyer(t)
		mockBuyer.On("Buy", "testUser", "t_shirt").Return(nil).Once()
		
		handler := buy.New(logger, mockBuyer)
		
		req := httptest.NewRequest(http.MethodPost, "/buy/t_shirt", nil)		
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("item", "t_shirt")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))
		
		w := httptest.NewRecorder()
		handler(w, req)
		
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("purchase error", func(t *testing.T) {		
		mockBuyer := mocks.NewBuyer(t)
		mockBuyer.On("Buy", "testUser", "t_shirt").Return(errors.New("purchase failed")).Once()

		handler := buy.New(logger, mockBuyer)

		req := httptest.NewRequest(http.MethodPost, "/buy/t_shirt", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("item", "t_shirt")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		userDTO := user.UserDTO{Username: "testUser"}
		req = req.WithContext(context.WithValue(req.Context(), user.UserDTOKey, userDTO))

		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing user info", func(t *testing.T) {		
		mockBuyer := mocks.NewBuyer(t)
		handler := buy.New(logger, mockBuyer)

		req := httptest.NewRequest(http.MethodPost, "/buy/t_shirt", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("item", "t_shirt")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))		

		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
