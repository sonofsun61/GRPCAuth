package hasher

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	saltStr := hex.EncodeToString(salt)
	return saltStr, nil
}

func HashPassword(password, salt string) (string, error) {
	str := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password's hash: %w", err)
	}
	return string(hash), nil
}

func ComparePassword(hash, password, salt string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
}
