package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	API
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	return &Config{
		API: NewAPI(),
	}
}
