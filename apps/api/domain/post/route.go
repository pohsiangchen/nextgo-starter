package post

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"apps/api/middleware"
)

func RegisterRoutes(r chi.Router, postCtrl *PostController, jwtMiddleware func(next http.Handler) http.Handler) {
	r.Route("/posts", func(r chi.Router) {
		r.Use(jwtMiddleware)

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.BindObject(&CreatePostRequest{}, postCtrl.validator))
			r.Post("/", postCtrl.Create)
		})

		r.Route("/{postID}", func(r chi.Router) {
			r.Use(SetPostToCtx(postCtrl.postService))
			r.Get("/", postCtrl.Get)

			r.Route("/", func(r chi.Router) {
				r.Use(middleware.BindObject(&UpdatePostRequest{}, postCtrl.validator))
				r.With(CheckPostOwnership).Patch("/", postCtrl.Update)
			})

			r.With(CheckPostOwnership).Delete("/", postCtrl.Delete)
		})
	})

	r.Route("/feeds", func(r chi.Router) {
		r.Use(jwtMiddleware)

		r.Route("/", func(r chi.Router) {
			r.Get("/", postCtrl.ListFeeds)
		})
	})
}
