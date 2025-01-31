package post

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"

	sqlcstore "apps/api/database/sqlc"
	"apps/api/middleware"
	"apps/api/util/response"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "post middleware context value " + k.name
}

var (
	// `PostCtxKey` is the context.Context key to store the `post` object
	PostCtxKey = &contextKey{"post"}
)

// Get the `post` object from `postID` of URL parameter and store it in the context.Context
func SetPostToCtx(postService PostService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			zlog := zerolog.Ctx(r.Context())

			idParam := chi.URLParam(r, "postID")
			postID, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				render.Render(w, r, response.ErrGeneric(http.StatusBadRequest, "Invalid URL parameter 'postID'"))
				return
			}

			ctx := r.Context()
			post, err := postService.FindById(ctx, postID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					render.Render(w, r, response.ErrNotFound)
				default:
					zlog.Error().Err(err).Msgf("failed to retrieve a post record with post ID '%d'", postID)
					render.Render(w, r, response.ErrInternalServerError)
				}
				return
			}

			// assigned `post` object to request context
			ctx = context.WithValue(ctx, PostCtxKey, post)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Get converted `postID` of URL parameter from the context.Context
func PostFromCtx(ctx context.Context) sqlcstore.Post {
	return ctx.Value(PostCtxKey).(sqlcstore.Post)
}

// Checks if post can be accessed by requesting user
func CheckPostOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		user, ok := middleware.UserFromCtx(ctx)
		if !ok {
			render.Render(w, r, response.ErrUnauthorized)
			return
		}
		post := PostFromCtx(ctx)

		if post.UserID != user.ID {
			render.Render(w, r, response.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
