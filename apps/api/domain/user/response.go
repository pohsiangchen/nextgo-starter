package user

import (
	"net/http"

	"github.com/go-chi/render"

	"apps/api/database/sqlc"
)

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func (hr UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewUserResponse(u sqlcstore.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}

func NewUserListResponse(users []sqlcstore.User) []render.Renderer {
	list := []render.Renderer{}
	for _, user := range users {
		ur := NewUserResponse(user)
		list = append(list, ur)
	}
	return list
}
