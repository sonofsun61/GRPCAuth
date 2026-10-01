package dto

type UserRegisterRequest struct {
	Email       string
	Name        string
	Surname     string
	PhoneNumber string
	Password    string
}

type UserLoginRequest struct {
	Email    string `db:"email"`
	Password string `db:"password_hash"`
}

type UserChangePasswordRequest struct {
	Email string 
	OldPassword string
	NewPassword string
}
