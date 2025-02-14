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
	Close() error
	TaskService
	// Add the new repository method here
}

// service estructura que contiene el cliente Mongo
type service struct {
	db *mongo.Client
	taskService
	// Add the new repository here
}

// New crea una nueva instancia de servicio de MongoDB
func New() Service {
	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	if host == "" || port == "" {
		log.Fatal("BLUEPRINT_DB_HOST or BLUEPRINT_DB_PORT environment variables are not set")
	}

	// Crear cliente de MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return &service{db: client}
}

// Health verifica el estado de la conexión con la base de datos
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.db.Ping(ctx, nil); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}

	return map[string]string{"status": "up", "message": "It's healthy"}
}

// Close cierra la conexión de MongoDB
func (s *service) Close() error {
	if err := s.db.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
		return err
	}
	return nil
}
