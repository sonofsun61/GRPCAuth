package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnString      string
	JWTSecret       string
	KafkaBrokers    []string
	UserEventsTopic string
}

func MustLoadConfig() *Config {
	_ = godotenv.Load()
	connString := os.Getenv("POSTGRES_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	userEventsTopic := os.Getenv("KAFKA_TOPIC_USER_EVENTS")
	if connString == "" {
		panic("POSTGRES_URL is not set")
	}
	if jwtSecret == "" {
		panic("JWT_SECRET is not set")
	}
	if kafkaBrokers == "" {
		panic("KAFKA_BROKERS is not set")
	}
	if userEventsTopic == "" {
		panic("KAFKA_TOPIC_USER_EVENTS is not set")
	}
	return &Config{
		ConnString:      connString,
		JWTSecret:       jwtSecret,
		KafkaBrokers:    strings.Split(kafkaBrokers, ","),
		UserEventsTopic: userEventsTopic,
	}
}
