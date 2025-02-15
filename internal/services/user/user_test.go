package user_test

import (
	"errors"
	"testing"
	"time"

	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/justcgh9/merch_store/internal/services"
	users "github.com/justcgh9/merch_store/internal/services/user"
	"github.com/justcgh9/merch_store/internal/services/user/mocks"
	"github.com/justcgh9/merch_store/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Authorize(t *testing.T) {
	accessSecret := "testsecret"
	mockRepo := mocks.NewUserRepo(t)
	t.Run("Successful authorization", func(t *testing.T) {
		mockRepo = mocks.NewUserRepo(t)
		psswd, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		mockRepo.On("GetUser", "testuser").Return(user.User{Username: "testuser", Password: string(psswd)}, nil)
		service := users.New(slog.Default(), accessSecret, mockRepo)

		token, err := service.Authorize("testuser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("User does not exist - Successful creation", func(t *testing.T) {
		mockRepo = mocks.NewUserRepo(t)
		mockRepo.On("GetUser", "newuser").Return(user.User{}, storage.ErrUserDoesNotExist)
		mockRepo.On("CreateUser", mock.Anything).Return(nil)
		service := users.New(slog.Default(), accessSecret, mockRepo)

		token, err := service.Authorize("newuser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		mockRepo = mocks.NewUserRepo(t)
		mockRepo.On("GetUser", "testuser").Return(user.User{Username: "testuser", Password: "$2a$10$5KPE6..."}, nil)
		service := users.New(slog.Default(), accessSecret, mockRepo)

		_, err := service.Authorize("testuser", "wrongpassword")
		assert.ErrorIs(t, err, services.UserIncorrectPassword)
	})

	t.Run("User creation fails", func(t *testing.T) {
		mockRepo = mocks.NewUserRepo(t)
		mockRepo.On("GetUser", "failuser").Return(user.User{}, storage.ErrUserDoesNotExist)
		mockRepo.On("CreateUser", mock.Anything).Return(errors.New("create error"))
		service := users.New(slog.Default(), accessSecret, mockRepo)

		_, err := service.Authorize("failuser", "password")
		assert.ErrorIs(t, err, services.UserRegistrationError)
	})
}

func TestUserService_Authenticate(t *testing.T) {
	accessSecret := "testsecret"
	mockRepo := mocks.NewUserRepo(t)
	service := users.New(slog.Default(), accessSecret, mockRepo)

	validToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"payload": user.UserDTO{
			Username: "validuser",
		},
	}).SignedString([]byte(accessSecret))

	invalidToken := "invalid.token.string"

	t.Run("Valid token", func(t *testing.T) {
		userDTO, err := service.Authenticate(validToken)
		assert.NoError(t, err)
		assert.Equal(t, "validuser", userDTO.Username)
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, err := service.Authenticate(invalidToken)
		assert.ErrorIs(t, err, services.UserErrInvalidToken)
	})
} 