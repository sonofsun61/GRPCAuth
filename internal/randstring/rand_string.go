package randstring

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var ErrInvalidLength error = errors.New("lenght must be more than 0")

func GenerateRandomPassword(lenght int) (string, error) {
	if lenght <= 0 {
		return "", ErrInvalidLength
	}
	newPassword := make([]byte, lenght)
	_, _ = rand.Read(newPassword)
	newPasswordStr := hex.EncodeToString(newPassword)
	return newPasswordStr, nil
}
