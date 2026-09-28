package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnString string
}

func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("Could not load .env file")
	}
	return &Config{
		ConnString: os.Getenv("DATABASE_URL"),
	}
}
