package services

import "errors"

var (
	UserRegistrationError    = errors.New("error creating new user")
	UserReadingError         = errors.New("error getting information about a user")
	UserIncorrectPassword    = errors.New("incorrect username or password")
	UserTokenGenerationError = errors.New("error generating token")
	UserErrInvalidToken      = errors.New("error invalid token")
	TransferZeroMoneyError   = errors.New("error cannot send less than 0 to another user")
	NonExistingItemError     = errors.New("given item does not exist")
	UnsuccessfulBuyError     = errors.New("buy operation did not succeed")
	GetInventoryError        = errors.New("could not get inventory")
	GetBalanceError          = errors.New("error accesing balance")
	GetHistoryError          = errors.New("error getting history")
)
