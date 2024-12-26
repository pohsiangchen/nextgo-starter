package middleware

import (
	"context"
	"net/http"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

// `RequestID` adapted from `RequestIDHandler` in https://github.com/rs/zerolog/blob/master/hlog/hlog.go

type idKey struct{}

// IDFromRequest returns the unique id associated to the request if any.
func IDFromRequest(r *http.Request) (id xid.ID, ok bool) {
	if r == nil {
		return
	}
	return IDFromCtx(r.Context())
}

// IDFromCtx returns the unique id associated to the context if any.
func IDFromCtx(ctx context.Context) (id xid.ID, ok bool) {
	id, ok = ctx.Value(idKey{}).(xid.ID)
	return
}

// CtxWithID adds the given xid.ID to the context
func CtxWithID(ctx context.Context, id xid.ID) context.Context {
	return context.WithValue(ctx, idKey{}, id)
}

// RequestIDHandler returns a handler setting a unique id to the request which can
// be gathered using IDFromRequest(req). This generated id is added as a field to the
// logger using the passed fieldKey as field name. The id is also added as a response
// header if the headerName is not empty.
//
// The generated id is a URL safe base64 encoded mongo object-id-like unique id.
// Mongo unique id generation algorithm has been selected as a trade-off between
// size and ease of use: UUID is less space efficient and snowflake requires machine
// configuration.
func RequestID(fieldKey, headerName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id, ok := IDFromRequest(r)
			if !ok {
				id = xid.New()
				ctx = CtxWithID(ctx, id)
				r = r.WithContext(ctx)
			}
			if fieldKey != "" {
				zlog := zerolog.Ctx(ctx)
				zlog.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str(fieldKey, id.String())
				})
			}
			if headerName != "" {
				w.Header().Set(headerName, id.String())
			}
			next.ServeHTTP(w, r)
		})
	}
}
