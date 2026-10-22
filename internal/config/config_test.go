package config_test

import (
	"authservice/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustLoadConfig_Success(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("JWT_SECRET", "test-secret")

	cfg := config.MustLoadConfig()
	require.NotNil(t, cfg)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.ConnString)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
}

func TestMustLoadConfig_MissingEnv(t *testing.T) {
	tests := []struct {
		name      string
		dbURL     string
		jwtSecret string
		wantPanic string
	}{
		{"no database url", "", "test-secret", "DATABASE_URL is not set"},
		{"no jwt secret", "postgres://user:pass@localhost:5432/db", "", "JWT_SECRET is not set"},
		{"nothing set", "", "", "DATABASE_URL is not set"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DATABASE_URL", tt.dbURL)
			t.Setenv("JWT_SECRET", tt.jwtSecret)

			assert.PanicsWithValue(t, tt.wantPanic, func() {
				config.MustLoadConfig()
			})
		})
	}
}
