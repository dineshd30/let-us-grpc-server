package main

import (
	"log"
	"net"

	"github.com/dineshd30/let-us-grpc-proto/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.GreeterServiceServer
}

const port = ":8080"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServiceServer(s, &server{})
	log.Printf("gRPC greet server listening on port %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
