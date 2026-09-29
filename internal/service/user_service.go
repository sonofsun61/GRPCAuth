package service

import (
	"authservice/internal/dto"
	"authservice/internal/entity"
	"authservice/internal/hasher"
	"context"
	"errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) error
}

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
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
	if err = s.repo.CreateUser(ctx, req); err != nil {
		return "", err
	}
	return "not implemented", nil
}
