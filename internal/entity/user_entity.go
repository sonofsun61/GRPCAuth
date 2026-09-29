package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	PhoneNumber  string    `db:"phone_number"`
	PasswordHash string    `db:"password_hash"`
	PasswordSalt string    `db:"password_salt"`
	CreatedAt     time.Time `db:"created_at"`
}
