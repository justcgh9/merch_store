package user

type UserDTO struct {
	Username string
}


func NewUserDTO(username string) *UserDTO {
	return &UserDTO{
		Username: username,
	}
}