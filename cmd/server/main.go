// package main contains the starting point of grpc Chat server
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mihailtudos/microservices/chat/internal/chat"
	desc "github.com/mihailtudos/microservices/chat/pkg/chat_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
	db      *pgxpool.Pool
	queries *chat.Queries
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	chatID := uuid.New()

	log.Printf("chat %s added users: %#v\n", chatID, req.GetUsernames())

	chat, err := s.queries.CreateChat(ctx, chat.CreateChatParams{
		Name:      gofakeit.DomainName(),
		Usernames: req.GetUsernames(),
	})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"failed to create new chat",
		)
	}

	log.Printf("chat %s created with %d users\n", chat.Name, len(chat.Usernames))

	return &desc.CreateResponse{
		Id: chatID.String(),
	}, nil
}

func setupDB(ctx context.Context) (*pgxpool.Pool, error) {
	dbString := os.Getenv("DATABASE_URL")
	if dbString == "" {
		dbString = os.Getenv("GOOSE_DBSTRING")
	}

	cfg, err := pgxpool.ParseConfig(dbString)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	// Sensible defaults (tune for prod)
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

func main() {
	ctx := context.Background()

	db, err := setupDB(ctx)
	if err != nil {
		log.Fatalf("failed to setupDB: %s", err)
	}

	queries := chat.New(db)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	server := &server{
		db:      db,
		queries: queries,
	}

	desc.RegisterChatV1Server(s, server)

	log.Printf("server listening at %s", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
