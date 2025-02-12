package storage

import "errors"

var (
	ErrUserDoesNotExist = errors.New("user with this username does not exist")
)
