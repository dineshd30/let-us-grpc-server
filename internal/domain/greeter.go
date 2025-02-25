package domain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/dineshd30/let-us-grpc-proto/proto"
)

type Server struct {
	proto.UnimplementedGreeterServiceServer
}

func (s *Server) SayHelloUnary(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	fmt.Println("---------------------------------------")
	log.Printf("Called SayHelloUnary: %s\n", req.Message)

	var response string
	if CheckVillein(req.Message) {
		response = fmt.Sprintf("Bye %s", req.Message)
	} else {
		response = fmt.Sprintf("Hello %s", req.Message)
	}

	log.Println("sending response:", response)
	return &proto.HelloResponse{
		Message: response,
	}, nil
}

func (s *Server) SayHelloServerStreaming(req *proto.NamesList, stream proto.GreeterService_SayHelloServerStreamingServer) error {
	fmt.Println("---------------------------------------")
	log.Println("Called SayHelloServerStreaming")

	for _, name := range req.Name {
		var response string
		if CheckVillein(name) {
			response = fmt.Sprintf("Bye %s", name)
		} else {
			response = fmt.Sprintf("Hello %s", name)
		}
		log.Println("sending response:", response)
		res := &proto.HelloResponse{
			Message: response,
		}

		if err := stream.Send(res); err != nil {
			return err
		}

		time.Sleep(2 * time.Second) // Artificial delay due to workload
	}

	return nil
}

func (s *Server) SayHelloClientStreaming(stream proto.GreeterService_SayHelloClientStreamingServer) error {
	fmt.Println("---------------------------------------")
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

		var response string
		if CheckVillein(req.Message) {
			response = fmt.Sprintf("Bye %s", req.Message)
		} else {
			response = fmt.Sprintf("Hello %s", req.Message)
		}
		log.Println("sending response:", response)
		messages = append(messages, response)
	}
}

func (s *Server) SayHelloBidirectionalStreaming(stream proto.GreeterService_SayHelloBidirectionalStreamingServer) error {
	fmt.Println("---------------------------------------")
	log.Println("Called SayHelloBidirectionalStreaming")

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			log.Fatalf("failed to get message from client stream: %v", err)
		}

		var response string
		if CheckVillein(req.Message) {
			response = fmt.Sprintf("Bye %s", req.Message)
		} else {
			response = fmt.Sprintf("Hello %s", req.Message)
		}
		log.Println("sending response:", response)
		res := &proto.HelloResponse{
			Message: response,
		}
		stream.Send(res)
		time.Sleep(time.Second * 2)
	}
}

func CheckVillein(name string) bool {
	villains := []string{"thanos", "ultron", "loki", "galactus", "kang the conqueror", "doctor doom"}
	return slices.IndexFunc(villains, func(c string) bool {
		return strings.ToLower(name) == c
	}) > -1
}
