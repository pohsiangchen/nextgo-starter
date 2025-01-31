package middleware

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Recovery adapted from https://github.com/go-chi/chi/blob/master/middleware/recoverer.go
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
					// HTTP server doesn't log an error, panic with the value `ErrAbortHandler` in Golang.
					// see https://pkg.go.dev/net/http
				}

				// converts panic to error
				var err error
				switch v := rvr.(type) {
				case string:
					err = errors.New(v)
				case error:
					err = v
				default:
					err = errors.New(fmt.Sprint(v))
				}

				zlog := zerolog.Ctx(r.Context())
				zlog.Error().Stack().Err(errors.WithStack(err)).Send()

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
