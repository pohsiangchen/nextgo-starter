package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"apps/api/util/response"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "object middleware context value " + k.name
}

var (
	ObjCtxKey = &contextKey{"obj"}
)

func Object(ctx context.Context) any {
	return ctx.Value(ObjCtxKey)
}

// Bind the request payload to target `obj` and perform validation if `validator` is not null
func BindObject(obj render.Binder, validator *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			zlog := zerolog.Ctx(r.Context())
			if err := render.Bind(r, obj); err != nil {
				zlog.Warn().Err(err).Msg("failed to unmarshal request payload")
				render.Render(w, r, response.ErrGeneric(http.StatusBadRequest, "Invalid request payload"))
				return
			}
			ctx := context.WithValue(r.Context(), ObjCtxKey, obj)

			if validator != nil {
				// validate incoming payload
				if err := validator.Struct(obj); err != nil {
					render.Render(w, r, response.ErrValidationFailed(err))
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
