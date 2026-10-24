package user

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrNameRequired = errors.New("name is mandatory")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}
