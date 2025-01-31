package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"

	"apps/api/database/sqlc"
	"apps/api/util/auth"
	"apps/api/util/response"
)

type jwtUserKey struct{}

// Get `user` object from associated context
func UserFromCtx(ctx context.Context) (sqlcstore.User, bool) {
	user, ok := ctx.Value(jwtUserKey{}).(sqlcstore.User)
	return user, ok
}

// Adds the given `user` object to the context
func CtxWithUser(ctx context.Context, user sqlcstore.User) context.Context {
	return context.WithValue(ctx, jwtUserKey{}, user)
}

func JWT(store *sqlcstore.Queries, authenticator auth.Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			zlog := zerolog.Ctx(r.Context())

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Render(w, r, response.ErrGeneric(http.StatusUnauthorized, "Authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				render.Render(w, r, response.ErrGeneric(http.StatusUnauthorized, "Authorization header is malformed"))
				return
			}

			token := parts[1]
			jwtToken, err := authenticator.ValidateToken(token)
			if err != nil {
				render.Render(w, r, response.ErrUnauthorized)
				return
			}

			claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
			if !ok {
				zlog.Warn().Msg("cannot parse claims from JWT token")
				render.Render(w, r, response.ErrUnauthorized)
				return
			}

			sub, err := claims.GetSubject()
			if err != nil {
				zlog.Warn().Msg("failed to get 'subject' from JWT claims")
				render.Render(w, r, response.ErrUnauthorized)
				return
			}

			userID, err := strconv.ParseInt(sub, 10, 64)
			if err != nil {
				render.Render(w, r, response.ErrUnauthorized)
				return
			}

			ctx := r.Context()

			user, err := store.GetUser(ctx, userID)
			if err != nil {
				render.Render(w, r, response.ErrUnauthorized)
				return
			}

			ctx = CtxWithUser(ctx, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
