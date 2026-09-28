package main

import (
	"authservice/gen/authpb"
	"authservice/internal/grpcserver"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to open tcp connection: %v", err)
	}
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, &grpcserver.Server{})
	reflection.Register(server)
	log.Println("Server started")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}