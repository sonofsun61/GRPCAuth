package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnString string
	JWTSecret  string
}

func MustLoadConfig() *Config {
	_ = godotenv.Load()
	connString := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if connString == "" {
		panic("DATABASE_URL is not set")
	}
	if jwtSecret == "" {
		panic("JWT_SECRET is not set")
	}
	return &Config{
		ConnString: connString,
		JWTSecret: jwtSecret,
	}
}
