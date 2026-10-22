package grpcserver_test

import (
	"authservice/internal/dto"
	"authservice/internal/grpcserver"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/sonofsun61/AuthContract/gen/authpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockAuthService struct {
	mock.Mock
}

var _ grpcserver.AuthService = (*mockAuthService)(nil)

func (m *mockAuthService) Register(ctx context.Context, newUserData dto.UserRegisterRequest) (string, error) {
	args := m.Called(ctx, newUserData)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, req dto.UserLoginRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockAuthService) ChangePassword(ctx context.Context, req dto.UserChangePasswordRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockAuthService) ResetPassword(ctx context.Context, req dto.UserResetPasswordRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

var errService = errors.New("service failed")

var discardLogger = slog.New(slog.DiscardHandler)

func requireGRPCCode(t *testing.T, err error, want codes.Code) {
	t.Helper()
	st, ok := status.FromError(err)
	require.True(t, ok, "error must be a gRPC status error")
	assert.Equal(t, want, st.Code())
}

func TestRegister_Success(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("Register", mock.Anything, dto.UserRegisterRequest{
		Email:       "example@mail.com",
		Name:        "Name",
		Surname:     "Surname",
		PhoneNumber: "+79999999999",
		Password:    "password",
	}).Return("token-123", nil)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.Register(context.Background(), &authpb.RegisterRequest{
		Email:     "example@mail.com",
		FirstName: "Name",
		LastName:  "Surname",
		Phone:     "+79999999999",
		Password:  "password",
	})
	require.NoError(t, err)
	assert.Equal(t, "token-123", resp.Token)
	svc.AssertExpectations(t)
}

func TestRegister_ServiceError(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("Register", mock.Anything, mock.Anything).Return("", errService)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.Register(context.Background(), &authpb.RegisterRequest{})
	requireGRPCCode(t, err, codes.Internal)
	assert.Nil(t, resp)
	svc.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("Login", mock.Anything, dto.UserLoginRequest{
		Email:    "example@mail.com",
		Password: "password",
	}).Return("token-123", nil)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.Login(context.Background(), &authpb.LoginRequest{
		Login:    "example@mail.com",
		Password: "password",
	})
	require.NoError(t, err)
	assert.Equal(t, "token-123", resp.Token)
	svc.AssertExpectations(t)
}

func TestLogin_ServiceError(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("Login", mock.Anything, mock.Anything).Return("", errService)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.Login(context.Background(), &authpb.LoginRequest{})
	requireGRPCCode(t, err, codes.Unauthenticated)
	assert.Nil(t, resp)
	svc.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("ChangePassword", mock.Anything, dto.UserChangePasswordRequest{
		Email:       "example@mail.com",
		OldPassword: "old-password",
		NewPassword: "new-password",
	}).Return(nil)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.ChangePassword(context.Background(), &authpb.ChangePasswordRequest{
		Login:       "example@mail.com",
		OldPassword: "old-password",
		NewPassword: "new-password",
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
	svc.AssertExpectations(t)
}

func TestChangePassword_ServiceError(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("ChangePassword", mock.Anything, mock.Anything).Return(errService)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.ChangePassword(context.Background(), &authpb.ChangePasswordRequest{})
	requireGRPCCode(t, err, codes.Unauthenticated)
	assert.Nil(t, resp)
	svc.AssertExpectations(t)
}

func TestResetPassword_Success(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("ResetPassword", mock.Anything, dto.UserResetPasswordRequest{
		Email: "example@mail.com",
	}).Return(nil)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.ResetPassword(context.Background(), &authpb.ResetPasswordRequest{
		Email: "example@mail.com",
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
	svc.AssertExpectations(t)
}

func TestResetPassword_ServiceError(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("ResetPassword", mock.Anything, mock.Anything).Return(errService)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.ResetPassword(context.Background(), &authpb.ResetPasswordRequest{})
	requireGRPCCode(t, err, codes.Internal)
	assert.Nil(t, resp)
	svc.AssertExpectations(t)
}

func TestValidateToken_Valid(t *testing.T) {
	svc := new(mockAuthService)
	userID := uuid.New()
	svc.On("ValidateToken", mock.Anything, "token-123").Return(userID, nil)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
		Token: "token-123",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, userID.String(), resp.UserId)
	svc.AssertExpectations(t)
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := new(mockAuthService)
	svc.On("ValidateToken", mock.Anything, "bad-token").Return(uuid.Nil, errService)
	server := grpcserver.NewServer(svc, discardLogger)

	resp, err := server.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
		Token: "bad-token",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Empty(t, resp.UserId)
	svc.AssertExpectations(t)
}
