package service_test

import (
	"authservice/internal/dto"
	"authservice/internal/entity"
	"authservice/internal/hasher"
	"authservice/internal/jwtutil"
	"authservice/internal/service"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	mock.Mock
}

var _ service.UserRepository = (*mockUserRepo)(nil)

func (m *mockUserRepo) CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, email, newPasswordHash, newPasswordSalt string) error {
	args := m.Called(ctx, email, newPasswordHash, newPasswordSalt)
	return args.Error(0)
}

const (
	testSecret = "test-secret-bot-for-production"
)

var errDBDown = errors.New("db down")

func validRegisterRequest() dto.UserRegisterRequest {
	req := dto.UserRegisterRequest{
		Email:       "example@gmail.com",
		Name:        "Name",
		Surname:     "Surname",
		PhoneNumber: "+79999999999",
		Password:    "password",
	}
	return req
}

func TestRegister_RepoError(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("CreateUser", mock.Anything, mock.Anything).Return(uuid.Nil, errDBDown)
	svc := service.NewAuthService(repo, testSecret)
	req := validRegisterRequest()

	token, err := svc.Register(context.Background(), req)
	assert.ErrorIs(t, err, errDBDown)
	assert.Empty(t, token)
	repo.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	repo := new(mockUserRepo)
	id := uuid.New()
	req := validRegisterRequest()
	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u entity.User) bool {
		return u.Email == req.Email && u.Name == req.Name &&
			u.Surname == req.Surname &&
			u.PhoneNumber == req.PhoneNumber &&
			u.PasswordHash != req.Password &&
			hasher.ComparePassword(u.PasswordHash, req.Password, u.PasswordSalt) == nil
	})).Return(id, nil)
	svc := service.NewAuthService(repo, testSecret)

	token, err := svc.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	userID, err := jwtutil.ValidateToken(token, testSecret)
	require.NoError(t, err)
	require.Equal(t, id, userID)
	repo.AssertExpectations(t)
}

func TestRegister_EmptyFields(t *testing.T) {
	tests := []struct {
		name   string
		modify func(r *dto.UserRegisterRequest)
	}{
		{"empty email", func(r *dto.UserRegisterRequest) { r.Email = "" }},
		{"empty name", func(r *dto.UserRegisterRequest) { r.Name = "" }},
		{"empty surname", func(r *dto.UserRegisterRequest) { r.Surname = "" }},
		{"empty phone number", func(r *dto.UserRegisterRequest) { r.PhoneNumber = "" }},
		{"empty password", func(r *dto.UserRegisterRequest) { r.Password = "" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockUserRepo)
			svc := service.NewAuthService(repo, testSecret)
			req := validRegisterRequest()
			tt.modify(&req)
			token, err := svc.Register(context.Background(), req)
			assert.Error(t, err)
			assert.Empty(t, token)
			repo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
		})
	}
}
