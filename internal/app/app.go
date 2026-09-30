package app

import (
	"authservice/gen/authpb"
	"authservice/internal/config"
	"authservice/internal/database"
	"authservice/internal/grpcserver"
	"authservice/internal/repository"
	"authservice/internal/service"
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	pool       *pgxpool.Pool
	grpcServer *grpc.Server
}


func NewApp() *App {
	cfg := config.MustLoadConfig()
	pool := database.MustConnectToDatabase(context.Background(), cfg.ConnString)

	repo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(repo)
	authServer := grpcserver.NewServer(authService)

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	return &App{
		pool:       pool,
		grpcServer: grpcServer,
	}
}

func (a *App) Run() {
	defer a.pool.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to open tcp connection: %v", err)
	}

	log.Printf("Server started on %s", lis.Addr())
	if err := a.grpcServer.Serve(lis); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}