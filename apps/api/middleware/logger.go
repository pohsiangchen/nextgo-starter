package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"apps/api/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		zlog := logger.Get()

		nwr := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		// The `Logger` instance is attached to Go context (context.Context) using `Logger.WithContext(ctx)`
		// and extracted from it using `zerolog.Ctx(ctx)`. We can log the additional fields from attached `context.Context`.
		r = r.WithContext(zlog.WithContext(r.Context()))

		defer func() {
			zerolog.Ctx(r.Context()).
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Str("ip", r.RemoteAddr).
				Int("status_code", nwr.Status()).
				Dur("elapsed_ms", time.Since(start)).
				Send()
		}()

		next.ServeHTTP(nwr, r)
	})
}
