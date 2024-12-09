package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (h HealthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{Status: "ok"}
	render.Render(w, r, health)
}