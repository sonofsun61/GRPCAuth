package entity

type User struct {
	Email        string `db:"email"`
	Name         string `db:"name"`
	Surname      string `db:"surname"`
	PhoneNumber  string `db:"phone_number"`
	PasswordHash string `db:"password_hash"`
	PasswordSalt string `db:"password_salt"`
}
