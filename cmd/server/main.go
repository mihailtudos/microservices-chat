// package main contains the starting point of grpc Chat server
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/mihailtudos/microservices/chat/pkg/chat_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	chatID := int64(12)
	log.Printf("chat %d added users: %#v\n\n", chatID, req.GetUsernames())

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %s", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
