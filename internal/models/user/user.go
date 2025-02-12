package user

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
