//go:build wireinject
// +build wireinject

package app

import (
	"authservice/internal/config"
	"authservice/internal/database"
	"authservice/internal/grpcserver"
	"authservice/internal/repository"
	"authservice/internal/service"
	"context"

	"github.com/sonofsun61/AuthContract/gen/authpb"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func InitializeApp() *App {
	wire.Build(
		config.MustLoadConfig,
		ProvideConnString,
		ProvideJWTSecret,
		ProvideContext,
		database.MustConnectToDatabase,
		repository.NewUserRepository,
		wire.Bind(new(service.UserRepository), new(*repository.UserRepository)),
		service.NewAuthService,
		wire.Bind(new(grpcserver.AuthService), new(*service.AuthService)),
		grpcserver.NewServer,
		ProvideGRPCServer,
		ProvideApp,
	)
	return nil
}

func ProvideConnString(cfg *config.Config) database.ConnString {
	return database.ConnString(cfg.ConnString)
}

func ProvideJWTSecret(cfg *config.Config) service.JWTSecret {
	return service.JWTSecret(cfg.JWTSecret)
}

func ProvideContext() context.Context {
	return context.Background()
}

func ProvideGRPCServer(authServer *grpcserver.Server) *grpc.Server {
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)
	return grpcServer
}

func ProvideApp(pool *pgxpool.Pool, grpcServer *grpc.Server) *App {
	return &App{
		pool:       pool,
		grpcServer: grpcServer,
	}
}
