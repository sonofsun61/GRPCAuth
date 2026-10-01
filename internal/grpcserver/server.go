package grpcserver

import (
	"authservice/gen/authpb"
	"authservice/internal/dto"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(ctx context.Context, newUserData dto.UserRegisterRequest) (string, error)
	Login(ctx context.Context, req dto.UserLoginRequest) (string, error)
	ChangePassword(ctx context.Context, req dto.UserChangePasswordRequest) error
}

type Server struct {
	authpb.UnimplementedAuthServiceServer
	authService AuthService
}

func NewServer(authService AuthService) *Server {
	return &Server{
		authService: authService,
	}
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	newUserData := dto.UserRegisterRequest{
		Email:       req.Email,
		Name:        req.FirstName,
		Surname:     req.LastName,
		PhoneNumber: req.Phone,
		Password:    req.Password,
	}
	token, err := s.authService.Register(ctx, newUserData)
	if err != nil {
		log.Printf("register error: %v", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &authpb.RegisterResponse{Token: token}, nil
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	loginRequest := dto.UserLoginRequest{
		Email:    req.Login,
		Password: req.Password,
	}
	token, err := s.authService.Login(ctx, loginRequest)
	if err != nil {
		log.Printf("login error: %v", err)
		return nil, status.Error(codes.Unauthenticated, "failed to login")
	}
	return &authpb.LoginResponse{Token: token}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	changePasswordRequest := dto.UserChangePasswordRequest{
		Email:       req.Login,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	if err := s.authService.ChangePassword(ctx, changePasswordRequest); err != nil {
		log.Printf("change password error: %v", err)
		return nil, status.Error(codes.Unauthenticated, "failed to update password")
	}
	return &authpb.ChangePasswordResponse{}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
