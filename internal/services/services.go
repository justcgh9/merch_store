package services

import "errors"

var (
	UserRegistrationError    = errors.New("error creating new user")
	UserReadingError         = errors.New("error getting information about a user")
	UserIncorrectPassword    = errors.New("incorrect username or password")
	UserTokenGenerationError = errors.New("error generating token")
	UserErrInvalidToken      = errors.New("error invalid token")
)
