package hasher_test

import (
	"authservice/internal/hasher"
	"authservice/internal/randstring"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateSalt_Format(t *testing.T) {
	salt, err := hasher.GenerateSalt()
	require.NoError(t, err)
	assert.Len(t, salt, 32)
	assert.Regexp(t, "^[0-9a-f]{32}$", string(salt))
}

func TestGenerateSalt_IsUnique(t *testing.T) {
	salt1, err := hasher.GenerateSalt()
	require.NoError(t, err)
	salt2, err := hasher.GenerateSalt()
	require.NoError(t, err)
	assert.NotEqual(t, salt1, salt2)
}

func TestHashPassword(t *testing.T) {
	password, err := randstring.GenerateRandomPassword(20)
	require.NoError(t, err)
	salt, err := hasher.GenerateSalt()
	require.NoError(t, err)
	hash, err := hasher.HashPassword(password, salt)
	require.NoError(t, err)
	require.NotEqual(t, hash, password+salt)
	assert.Regexp(t, `^\$2`, hash)
	hash2, err := hasher.HashPassword(password, salt)
	require.NoError(t, err)
	require.NotEqual(t, hash, hash2)
}

func TestComparePassword(t *testing.T) {
	password, err := randstring.GenerateRandomPassword(20)
	require.NoError(t, err)
	salt, err := hasher.GenerateSalt()
	require.NoError(t, err)
	hash, err := hasher.HashPassword(password, salt)
	require.NoError(t, err)
	tests := []struct {
		name     string
		hash     string
		password string
		salt     string
		wantErr  error
	}{
		{"RightPasswordAndSalt", hash, password, salt, nil},
		{"WrongPassword", hash, "wrongpassword", salt, bcrypt.ErrMismatchedHashAndPassword},
		{"WrongSalt", hash, password, "wrongsalt", bcrypt.ErrMismatchedHashAndPassword},
		{"EmptyPassword", hash, "", salt, bcrypt.ErrMismatchedHashAndPassword},
		{"UpperCasePassword", hash, strings.ToUpper(password), salt, bcrypt.ErrMismatchedHashAndPassword},
		{"NotAHash", "fjwlnfvklgnr", password, salt, bcrypt.ErrHashTooShort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := hasher.ComparePassword(tt.hash, tt.password, tt.salt)
			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
			} else {
				require.NoError(t, gotErr)
			}
		})
	}
}
