package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"apps/api/config"
)

// Adapted from https://github.com/betterstack-community/wikipedia-demo/blob/zerolog/logger/logger.go

var once sync.Once
var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		cfg := config.Get()
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
			},
		}

		if cfg.API.AppEnv != "development" {
			output = os.Stdout
		}

		log = zerolog.New(output).
			Level(zerolog.Level(cfg.API.LogLevel)).
			With().
			Timestamp().
			Logger()

		zerolog.DefaultContextLogger = &log
	})

	return log
}
