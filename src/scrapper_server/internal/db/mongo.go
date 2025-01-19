package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	Close()
}

type service struct {
	db *mongo.Client
}

func New() Service {
	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	if host == "" || port == "" {
		log.Fatal("BLUEPRINT_DB_HOST or BLUEPRINT_DB_PORT environment variables are not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return &service{db: client}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.Ping(ctx, nil); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}

	return map[string]string{"status": "up", "message": "It's healthy"}
}

func (s *service) Close() {
	if err := s.db.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	}
}
