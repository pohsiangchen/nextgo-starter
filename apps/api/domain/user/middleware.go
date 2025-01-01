package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"apps/api/util/response"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "user middleware context value " + k.name
}

var (
	// `UserIDCtxKey` is the context.Context key to store the `userID` from URL parameter
	UserIDCtxKey = &contextKey{"userID"}
)

// Convert `userID` of URL parameter to `int64` and store it in the context.Context
func SetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			render.Render(w, r, response.ErrGeneric(http.StatusBadRequest, "Invalid URL parameter 'userID'"))
			return
		}
		// assigned converted `userID` to request context
		ctx := context.WithValue(r.Context(), UserIDCtxKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get converted `userID` of URL parameter from the context.Context
func UserID(ctx context.Context) int64 {
	return ctx.Value(UserIDCtxKey).(int64)
}
