package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"greet_server/greetpb"
	"greet_server/server"
)

func main() {
	fmt.Println("Hello from gRPC's Server")

	gRPCSever := server.NewServerInstance()

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// certFile := "../ssl/server.crt"
	// keyFile := "../ssl/server.key"
	// creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	// if sslErr != nil {
	// 	log.Fatalf("Error TLS: %v", sslErr)
	// 	return
	// }

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, gRPCSever)

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
