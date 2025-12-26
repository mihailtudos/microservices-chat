// package main contains the starting point of grpc Chat client
package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	desc "github.com/mihailtudos/microservices/chat/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
const (
	address = "localhost:50052"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed to connect to server: %s", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	c := desc.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var usernames []string

	for range 5 {
		usernames = append(usernames, gofakeit.Username())
	}

	r, err := c.Create(ctx, &desc.CreateRequest{Usernames: usernames})
	if err != nil {
		log.Fatalf("failed to create char for users: %#v due to: %s", usernames, err)
	}

	log.Printf("note info: %+v", r.GetId())
}
