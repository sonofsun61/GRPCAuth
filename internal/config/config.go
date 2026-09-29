package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnString string
}

func MustLoadConfig() *Config {
	_ = godotenv.Load()
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		panic("DATABASE_URL is not set")
	}
	return &Config{
		ConnString: connString,
	}
}
