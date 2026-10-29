package config_test

import (
	"authservice/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustLoadConfig_Success(t *testing.T) {
	t.Setenv("POSTGRES_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("KAFKA_BROKERS", "k1:9092,k2:9092")
	t.Setenv("KAFKA_TOPIC_USER_EVENTS", "user.events")

	cfg := config.MustLoadConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.ConnString)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
	assert.Equal(t, []string{"k1:9092", "k2:9092"}, cfg.KafkaBrokers)
	assert.Equal(t, "user.events", cfg.UserEventsTopic)
}

func TestMustLoadConfig_MissingEnv(t *testing.T) {
	valid := map[string]string{
		"POSTGRES_URL":            "postgres://user:pass@localhost:5432/db",
		"JWT_SECRET":              "test-secret",
		"KAFKA_BROKERS":           "localhost:9092",
		"KAFKA_TOPIC_USER_EVENTS": "user.events",
	}
	tests := []struct {
		name      string
		unset     string
		wantPanic string
	}{
		{"no database url", "POSTGRES_URL", "POSTGRES_URL is not set"},
		{"no jwt secret", "JWT_SECRET", "JWT_SECRET is not set"},
		{"no kafka brokers", "KAFKA_BROKERS", "KAFKA_BROKERS is not set"},
		{"no user events topic", "KAFKA_TOPIC_USER_EVENTS", "KAFKA_TOPIC_USER_EVENTS is not set"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range valid {
				t.Setenv(k, v)
			}
			t.Setenv(tt.unset, "")

			assert.PanicsWithValue(t, tt.wantPanic, func() {
				config.MustLoadConfig()
			})
		})
	}
}
