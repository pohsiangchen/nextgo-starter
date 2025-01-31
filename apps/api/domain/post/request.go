package post

import (
	"net/http"
)

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1"`
}

type UpdatePostRequest struct {
	ID      int64  `json:"-"`
	Title   string `json:"title" validate:"omitempty,min=1"`
	Content string `json:"content" validate:"omitempty,min=1"`
}

func (cur *CreatePostRequest) Bind(r *http.Request) error {
	return nil
}

func (uur *UpdatePostRequest) Bind(r *http.Request) error {
	return nil
}
