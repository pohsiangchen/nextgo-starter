package auth

import "net/http"

type CreateTokenRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=16"`
}

func (cur *CreateTokenRequest) Bind(r *http.Request) error {
	return nil
}
