package user

import "github.com/golang-jwt/jwt/v5"

type userDTOKey string

const UserDTOKey userDTOKey = "userDTO"

type UserDTO struct {
	Username string
}

func NewUserDTO(username string) *UserDTO {
	return &UserDTO{
		Username: username,
	}
}

type User struct {
	Username string
	Password string
}

type UserClaims struct {
	Payload UserDTO `json:"payload"`
	jwt.RegisteredClaims
}
