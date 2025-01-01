package user

import (
	"apps/api/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, userCtrl *UserController) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", userCtrl.List)

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.BindObject(&CreateUserRequest{}, userCtrl.validator))
			r.Post("/", userCtrl.Create)
		})

		r.Route("/{userID}", func(r chi.Router) {
			r.Use(SetUserID)
			r.Get("/", userCtrl.Get)

			r.Route("/", func(r chi.Router) {
				r.Use(middleware.BindObject(&UpdateUserRequest{}, userCtrl.validator))
				r.Put("/", userCtrl.Update)
			})

			r.Delete("/", userCtrl.Delete)
		})

	})
}
