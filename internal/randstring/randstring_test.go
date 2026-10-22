package randstring_test

import (
	"authservice/internal/randstring"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomPassword_LengthCheck(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"zero length", 0, true},
		{"negative length", -1, true},
		{"one length", 1, false},
		{"twelve length", 12, false},
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
				assert.Len(t, password, 2*tt.length)
				assert.Regexp(t, "^[0-9a-f]+$", password)
			}
		})
	}
}

func TestGenerateRandomPassword_IsUnique(t *testing.T) {
	password1, err := randstring.GenerateRandomPassword(12)
	require.NoError(t, err)
	password2, err := randstring.GenerateRandomPassword(12)
	require.NoError(t, err)
	assert.NotEqual(t, password1, password2)
}
