package database

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var (
	mongoTeardown func(context.Context) error
	mongoURI      string
)

// Configuración del contenedor MongoDB
func mustStartMongoContainer() (func(context.Context) error, string, error) {
	ctx := context.Background()
	// Crear y configurar contenedor de MongoDB
	dbContainer, err := mongodb.Run(ctx, "mongo:latest")
	if err != nil {
		return nil, "", err
	}

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return func(ctx context.Context) error {
			return dbContainer.Terminate(ctx)
		}, "", err
	}

	port, err := dbContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return func(ctx context.Context) error {
			return dbContainer.Terminate(ctx)
		}, "", err
	}

	uri := "mongodb://" + host + ":" + port.Port()
	return func(ctx context.Context) error {
		return dbContainer.Terminate(ctx)
	}, uri, nil
}

// Preparar entorno de pruebas
func TestMain(m *testing.M) {
	var err error
	mongoTeardown, mongoURI, err = mustStartMongoContainer()
	if err != nil {
		log.Fatalf("Could not start MongoDB container: %v", err)
	}
	log.Printf("MongoDB container running at: %s", mongoURI)

	// Configurar variables de entorno para el servicio
	os.Setenv("BLUEPRINT_DB_HOST", mongoURI)
	os.Setenv("BLUEPRINT_DB_PORT", "27017")

	// Ejecutar pruebas
	code := m.Run()

	// Finalizar contenedor
	if mongoTeardown != nil {
		if err := mongoTeardown(context.Background()); err != nil {
			log.Printf("Failed to teardown MongoDB container: %v", err)
		}
	}

	os.Exit(code)
}

// Test de creación del servicio
func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

// Test del método Health en caso exitoso
func TestHealth(t *testing.T) {
	srv := New()
	defer srv.Close()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("Expected status to be 'up', got %s", stats["status"])
	}
	if stats["message"] != "It's healthy" {
		t.Fatalf("Expected message to be 'It's healthy', got %s", stats["message"])
	}
}

// Test del método Health en caso fallido
func TestHealth_Down(t *testing.T) {
	// Crear un servicio y cerrar manualmente el cliente MongoDB para simular una caída
	srv := New()
	srv.Close()

	stats := srv.Health()

	if stats["status"] != "down" {
		t.Fatalf("Expected status to be 'down', got %s", stats["status"])
	}
	if stats["error"] == "" {
		t.Fatal("Expected an error message, but got none")
	}
}
