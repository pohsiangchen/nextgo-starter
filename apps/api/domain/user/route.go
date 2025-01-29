package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"apps/api/middleware"
)

func RegisterRoutes(r chi.Router, userCtrl *UserController, jwtMiddleware func(next http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Use(jwtMiddleware)

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
