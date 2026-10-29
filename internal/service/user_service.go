package service

import (
	"authservice/internal/dto"
	"authservice/internal/entity"
	"authservice/internal/events"
	"authservice/internal/hasher"
	"authservice/internal/jwtutil"
	"authservice/internal/randstring"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type JWTSecret string

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UpdatePassword(ctx context.Context, email, newPasswordHash, newPasswordSalt string) error
}

type EventPublisher interface {
	PublishUserRegistered(ctx context.Context, eventInfo events.UserRegistered) error
}

type AuthService struct {
	repo      UserRepository
	publisher EventPublisher
	jwtSecret string
	logger    *slog.Logger
}

func NewAuthService(repo UserRepository, publisher EventPublisher, jwtSecret JWTSecret, logger *slog.Logger) *AuthService {
	return &AuthService{
		repo:      repo,
		publisher: publisher,
		jwtSecret: string(jwtSecret),
		logger:    logger,
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
	eventInfo := events.UserRegistered{
		EventID:    uuid.New(),
		UserID:     id,
		Email:      req.Email,
		Name:       req.Name,
		OccurredAt: time.Now().UTC(),
	}
	if err := s.publisher.PublishUserRegistered(ctx, eventInfo); err != nil {
		s.logger.Warn("failed to publish create user event", "user_id", id, "error", err)
	}
	token, err := jwtutil.GenerateToken(id, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.UserLoginRequest) (string, error) {
	userInfo, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	err = hasher.ComparePassword(userInfo.PasswordHash, req.Password, userInfo.PasswordSalt)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	token, err := jwtutil.GenerateToken(userInfo.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req dto.UserChangePasswordRequest) error {
	userInfo, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("change password: user lookup failed", "email", req.Email, "err", err)
		return errors.New("invalid email or password")
	}
	err = hasher.ComparePassword(userInfo.PasswordHash, req.OldPassword, userInfo.PasswordSalt)
	if err != nil {
		return errors.New("invalid email or password")
	}
	newSalt, err := hasher.GenerateSalt()
	if err != nil {
		return err
	}
	newPassword, err := hasher.HashPassword(req.NewPassword, newSalt)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, req.Email, newPassword, newSalt); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req dto.UserResetPasswordRequest) error {
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("reset password: user lookup failed", "email", req.Email, "err", err)
		return nil
	}
	newPassword, err := randstring.GenerateRandomPassword(12)
	if err != nil {
		return err
	}
	salt, err := hasher.GenerateSalt()
	if err != nil {
		return err
	}
	hashPassword, err := hasher.HashPassword(newPassword, salt)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(ctx, req.Email, hashPassword, salt); err != nil {
		return err
	}
	// TODO: add email sending
	s.logger.Info("password reset", "email", req.Email, "new_password", newPassword)
	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	return jwtutil.ValidateToken(token, s.jwtSecret)
}
