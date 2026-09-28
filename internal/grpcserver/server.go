package grpcserver

import (
	"authservice/gen/authpb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *Server) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *Server) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
