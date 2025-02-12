package user

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justcgh9/merch_store/internal/models/user"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) createUser(username, password string) error {

	log := u.log.With(
		slog.String("username", username),
	)

	log.Info("creating user")

	pswd, err := hashPassword(password)
	if err != nil {
		log.Error("error hashing password", slog.String("err", err.Error()))
		return err
	}

	err = u.userRepo.CreateUser(user.User{
		Username: username,
		Password: pswd,
	})
	if err != nil {
		log.Error("error creating user", slog.String("err", err.Error()))
		return err
	}

	log.Info("created user")

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateTokens(accessSecret, username string) (string, error) {

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"payload": user.UserDTO{
			Username: username,
		},
	}).SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
