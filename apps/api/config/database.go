package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Driver                 string        `required:"true"`
	Host                   string        `default:"localhost"`
	Port                   uint16        `default:"5432"`
	Name                   string        `default:"nextgo_db"`
	User                   string        `default:"postgres"`
	Password               string        `default:"password"`
	SslMode                string        `split_words:"true" default:"disable"`
	MaxConnectionPool      int           `split_words:"true" default:"4"`
	MaxIdleConnections     int           `split_words:"true" default:"4"`
	ConnectionsMaxLifeTime time.Duration `split_words:"true" default:"300s"`
}

func NewDatabase() Database {
	var db Database
	envconfig.MustProcess("DB", &db)

	return db
}
