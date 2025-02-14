package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	database "scrapper_server/internal/db"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port       int
	db         database.Service
	workerPool workerPool
}

func NewServer(db database.Service) *http.Server {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT environment variable: %s", portStr)
	}

	srv := &Server{
		port:       port,
		db:         db,
		workerPool: NewWorkerPool(5, 100), // Crear un pool de 5 workers con capacidad de 100 tareas en cola
	}

	// Configuración del servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
