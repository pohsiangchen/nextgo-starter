package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"apps/api/middleware"
)

func RegisterRoutes(r chi.Router, userCtrl *UserController, jwtMiddleware func(next http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Use(jwtMiddleware)

		// TODO: add user role feature to allow admin role only
		// r.Get("/", userCtrl.List)

		// TODO: add user role feature to allow admin role only
		// r.Route("/", func(r chi.Router) {
		// 	r.Use(middleware.BindObject(&CreateUserRequest{}, userCtrl.validator))
		// 	r.Post("/", userCtrl.Create)
		// })

		r.Route("/{userID}", func(r chi.Router) {
			r.Use(SetUserIDToCtx)
			r.With(CheckUserOwnership).Get("/", userCtrl.Get)

			r.Route("/", func(r chi.Router) {
				r.Use(middleware.BindObject(&UpdateUserRequest{}, userCtrl.validator))
				r.With(CheckUserOwnership).Put("/", userCtrl.Update)
			})

			// TODO: add user role feature to allow admin role only
			// r.Delete("/", userCtrl.Delete)
		})

	})
}
