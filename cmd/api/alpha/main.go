package main

// import (
// 	"log"
// 	"net"

// 	"google.golang.org/grpc"

// 	"go.breu.io/quantm/internal/nomad/handler"
// 	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
// 	healthzv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/healthz/v1"
// )

// func main() {
// 	listen, err := net.Listen("tcp", ":7070") // nolint
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	srv := grpc.NewServer()

// 	// Register the service
// 	authv1.RegisterAccountServiceServer(srv, &handler.AccountService{})
// 	healthzv1.RegisterHealthCheckServiceServer(srv, &handler.HealthCheckService{})

// 	_ = srv.Serve(listen)
// }
