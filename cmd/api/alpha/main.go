package main

// import (
// 	"log"
// 	"net"

// 	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
// 	"google.golang.org/grpc"
// )

// func main() {
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()

// 	// Register the service
// 	authv1.RegisterAccountServiceServer(s, &AccountService{})
// 	s.Serve(lis)
// }
