package hasher_test

import (
	"authservice/internal/hasher"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSalt_Format(t *testing.T) {
	salt, err := hasher.GenerateSalt()
	require.NoError(t, err)
	assert.Len(t, salt, 32)
	assert.Regexp(t, regexp.MustCompile(string(salt)), string(salt))
}

func TestGenerateSalt_IsUnique(t *testing.T) {
	salt1, err := hasher.GenerateSalt()
	require.NoError(t, err)
	salt2, err := hasher.GenerateSalt()
	require.NoError(t, err)
	assert.NotEqual(t, salt1, salt2)
}
