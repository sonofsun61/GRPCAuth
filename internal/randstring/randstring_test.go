package randstring_test

import (
	"authservice/internal/randstring"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomPasword_LengthCheck(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"zero length", 0, true},
		{"negative length", -1, true},
		{"one length", 1, false},
		{"twelwe length", 12, false},
		{"twenty length", 20, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := randstring.GenerateRandomPassword(tt.length)
			if tt.wantErr {
				require.ErrorIs(t, err, randstring.ErrInvalidLength)
				assert.Equal(t, "", password)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
