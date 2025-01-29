package auth

import (
	"github.com/go-chi/chi/v5"

	"apps/api/middleware"
)

func RegisterRoutes(r chi.Router, authCtrl *AuthController) {
	r.Route("/authentication", func(r chi.Router) {

		r.Route("/token", func(r chi.Router) {
			r.Use(middleware.BindObject(&CreateTokenRequest{}, authCtrl.validator))
			r.Post("/", authCtrl.CreateToken)
		})
	})
}
