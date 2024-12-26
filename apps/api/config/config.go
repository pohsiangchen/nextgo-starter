package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	API
	Database
}

var once sync.Once
var cfg *Config

func Get() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println(err)
		}

		cfg = &Config{
			API:      NewAPI(),
			Database: NewDatabase(),
		}
	})

	return cfg
}
