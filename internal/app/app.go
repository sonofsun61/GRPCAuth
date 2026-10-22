package app

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	pool       *pgxpool.Pool
	grpcServer *grpc.Server
	logger     *slog.Logger
}

func (a *App) Run() error {
	defer a.pool.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to open tcp connection: %w", err)
	}

	a.logger.Info("server started", "addr", lis.Addr().String())
	if err := a.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("server failed to serve: %w", err)
	}
	return nil
}
