package app

import (
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	pool       *pgxpool.Pool
	grpcServer *grpc.Server
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
