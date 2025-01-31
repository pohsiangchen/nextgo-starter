package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type API struct {
	Name            string        `default:"nextgo-starter"`
	Host            string        `default:"localhost"`
	Port            string        `default:"8080"`
	IdleTimeout     time.Duration `split_words:"true" default:"60s"`
	ReadTimeout     time.Duration `split_words:"true" default:"5s"`
	WriteTimeout    time.Duration `split_words:"true" default:"10s"`
	GracefulTimeout time.Duration `split_words:"true" default:"10s"`
	AppEnv          string        `split_words:"true" default:"development"`
	LogLevel        int           `split_words:"true" default:"1"`
	AuthJwtSecret   string        `required:"true" split_words:"true"`
	AuthJwtIss      string        `required:"true" split_words:"true"`
	AuthJwtExp      time.Duration `split_words:"true" default:"72h"`
}

func NewAPI() API {
	var api API
	envconfig.MustProcess("API", &api)

	return api
}
