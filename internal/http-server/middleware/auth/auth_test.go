package auth_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justcgh9/merch_store/internal/http-server/middleware/auth"
	"github.com/justcgh9/merch_store/internal/http-server/middleware/auth/mocks"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	authenticatorMock := mocks.NewAuthenticator(t)
	middleware := auth.New(log, authenticatorMock)

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "missing authorization header")
	})

	t.Run("invalid authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()

		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid authorization header")
	})

	t.Run("authentication failure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		w := httptest.NewRecorder()

		authenticatorMock.On("Authenticate", "invalidtoken").Return(user.UserDTO{}, errors.New("invalid jwt token"))

		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid jwt token")
	})

	t.Run("authentication success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer validtoken")
		w := httptest.NewRecorder()

		authenticatorMock.On("Authenticate", "validtoken").Return(user.UserDTO{Username: "testuser"}, nil)

		middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userDTO, ok := r.Context().Value(user.UserDTOKey).(user.UserDTO)
			assert.True(t, ok)
			assert.Equal(t, "testuser", userDTO.Username)
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
