package randstring

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var ErrInvalidLength = errors.New("length must be more than 0")

func GenerateRandomPassword(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}
	newPassword := make([]byte, length)
	_, _ = rand.Read(newPassword)
	newPasswordStr := hex.EncodeToString(newPassword)
	return newPasswordStr, nil
}
