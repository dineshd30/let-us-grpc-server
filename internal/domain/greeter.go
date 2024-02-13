package domain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dineshd30/let-us-grpc-proto/proto"
)

type Server struct {
	proto.UnimplementedGreeterServiceServer
}

func (s *Server) SayHelloUniary(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Called SayHelloUniary: %s\n", req.Message)

	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.Message),
	}, nil
}

func (s *Server) SayHelloServerStreaming(req *proto.NamesList, stream proto.GreeterService_SayHelloServerStreamingServer) error {
	log.Println("Called SayHelloServerStreaming")

	for _, name := range req.Name {
		res := &proto.HelloResponse{
			Message: fmt.Sprintf("Hello %s", name),
		}

		if err := stream.Send(res); err != nil {
			return err
		}

		time.Sleep(2 * time.Second) // Artificial delay due to workload
	}

	return nil
}

func (s *Server) SayHelloClientStreaming(stream proto.GreeterService_SayHelloClientStreamingServer) error {
	log.Println("Called SayHelloClientStreaming")

	var messages []string
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&proto.MessageList{
				Name: messages,
			})
		}
		if err != nil {
			log.Fatalf("failed to get message from client stream: %v", err)
		}
		log.Printf("received message from client stream - %s", req.Message)
		messages = append(messages, fmt.Sprintf("Hello %s", req.Message))
	}
}

func (s *Server) SayHelloBidirectionalStreaming(stream proto.GreeterService_SayHelloBidirectionalStreamingServer) error {
	log.Println("Called SayHelloBidirectionalStreaming")

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			log.Fatalf("failed to get message from client stream: %v", err)
		}

		log.Printf("received message from client stream - %s", req.Message)
		res := &proto.HelloResponse{
			Message: fmt.Sprintf("Hello %s", req.Message),
		}
		stream.Send(res)
		log.Printf("Sent response to client - %s", res.Message)
		time.Sleep(time.Second * 2)
	}
}
