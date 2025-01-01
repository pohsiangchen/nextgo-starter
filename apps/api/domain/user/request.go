package user

import (
	"net/http"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=2,max=32"`
	Password string `json:"password" validate:"required,min=6,max=16"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"-"`
	Username string `json:"username" validate:"required,min=2,max=32"`
}

type UpdateUserPasswordRequest struct {
	ID       int64  `json:"-"`
	Password string `json:"password" validate:"required,min=6"`
}

func (cur *CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

func (uur *UpdateUserRequest) Bind(r *http.Request) error {
	return nil
}

// TDOD: update user password
func (uupr *UpdateUserPasswordRequest) Bind(r *http.Request) error {
	return nil
}
