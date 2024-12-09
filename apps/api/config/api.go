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
}

func NewAPI() API {
	var api API
	envconfig.MustProcess("API", &api)

	return api
}
