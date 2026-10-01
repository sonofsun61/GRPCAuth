package randstring

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomPassword(lenght int) (string, error) {
	newPassword := make([]byte, lenght)
	_, err := rand.Read(newPassword)
	if err != nil {
		return "", fmt.Errorf("failed to generate new password: %w", err)
	}
	newPasswordStr := hex.EncodeToString(newPassword)
	return newPasswordStr, nil
}