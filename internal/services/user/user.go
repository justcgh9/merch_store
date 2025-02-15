package user

import (
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/justcgh9/merch_store/internal/services"
	"github.com/justcgh9/merch_store/internal/storage"
)

type UserRepo interface {
	GetUser(username string) (user.User, error)
	CreateUser(user user.User) error
}

type UserService struct {
	log          *slog.Logger
	userRepo     UserRepo
	accessSecret string
}

func New(log *slog.Logger, accessSecret string, userRepo UserRepo) *UserService {
	return &UserService{
		log:          log,
		accessSecret: accessSecret,
		userRepo:     userRepo,
	}
}

func (u *UserService) Authenticate(tokenStr string) (user.UserDTO, error) {
	const op = "services.user.Authenticate"

	log := u.log.With(
		slog.String("op", op),
	)

	log.Info("authenticating user")

	secretKey := []byte(u.accessSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &user.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok && token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			return secretKey, nil
		}
		return nil, services.UserErrInvalidToken
	})

	if err != nil {
		log.Error("invalid jwt token", slog.String("err", err.Error()))
		return user.UserDTO{}, services.UserErrInvalidToken
	}

	if claims, ok := token.Claims.(*user.UserClaims); ok && token.Valid {
		log.Info("token validated successfully", slog.Any("username", claims.Payload.Username))
		return claims.Payload, nil
	}

	log.Error("invalid jwt token")
	return user.UserDTO{}, services.UserErrInvalidToken
}

func (u *UserService) Authorize(username, password string) (string, error) {
	const op = "services.user.Authorize"

	log := u.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("authorizing user")

	user, err := u.userRepo.GetUser(username)
	if err != nil {
		if errors.Is(err, storage.ErrUserDoesNotExist) {
			err := u.createUser(username, password)
			if err != nil {
				log.Error("error creating user", slog.String("err", err.Error()))
				return "", services.UserRegistrationError
			}
			user.Username = username
		} else {
			log.Error("error reading user", slog.String("err", err.Error()))
			return "", services.UserReadingError
		}
	} else {

		if !checkPasswordHash(password, user.Password) {
			log.Error("incorrect username or password")
			return "", services.UserIncorrectPassword
		}

	}

	token, err := generateTokens(u.accessSecret, user.Username)
	if err != nil {
		log.Error("error generating token", slog.String("err", err.Error()))
		return "", services.UserTokenGenerationError
	}

	return token, nil
}
