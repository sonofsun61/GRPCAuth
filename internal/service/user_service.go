package service

import (
	"authservice/internal/dto"
	"authservice/internal/entity"
	"authservice/internal/hasher"
	"authservice/internal/jwtutil"
	"context"
	"errors"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
}

type AuthService struct {
	repo      UserRepository
	jwtSecret string
}

func NewAuthService(repo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, newUserData dto.UserRegisterRequest) (string, error) {
	if len(newUserData.Email) == 0 || len(newUserData.Password) == 0 ||
		len(newUserData.Name) == 0 || len(newUserData.Surname) == 0 ||
		len(newUserData.PhoneNumber) == 0 {
		return "", errors.New("request data is incorrect")
	}
	salt, err := hasher.GenerateSalt()
	if err != nil {
		return "", err
	}
	passwordHash, err := hasher.HashPassword(newUserData.Password, salt)
	if err != nil {
		return "", err
	}
	req := entity.User{
		Email:        newUserData.Email,
		Name:         newUserData.Name,
		Surname:      newUserData.Surname,
		PhoneNumber:  newUserData.PhoneNumber,
		PasswordHash: passwordHash,
		PasswordSalt: salt,
	}
	id, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		return "", err
	}
	token, err := jwtutil.GenerateToken(id, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}
