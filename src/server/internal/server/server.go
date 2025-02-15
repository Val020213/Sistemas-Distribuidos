package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"server/internal/database"
	"server/internal/models"
	"server/internal/scraper"
	"server/internal/utils"
)

type Server struct {
	port      int
	db        database.Service
	TaskQueue chan string
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	queueSize := utils.GetEnvAsInt("TASK_QUEUE_SIZE", 100)
	numWorkers := utils.GetEnvAsInt("NUM_WORKERS", 5)

	NewServer := &Server{
		port:      port,
		db:        database.New(),
		TaskQueue: make(chan string, queueSize),
	}

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		go NewServer.worker()
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) worker() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker recovered from panic: %v", r)
		}
	}()

	for taskID := range s.TaskQueue {

		taskIDUint, err := strconv.ParseUint(taskID, 10, 32)
		if err != nil {
			log.Printf("Invalid task ID format %s: %v", taskID, err)
			continue
		}

		task, err := s.db.GetTask(uint32(taskIDUint))
		if err != nil {
			log.Printf("Error fetching task %d: %v", taskIDUint, err)
			continue
		}

		// Scrape URL
		content, err := scraper.ScrapeURL(task.URL)
		if err != nil {
			log.Printf("Error scraping URL %s: %v", task.URL, err)
			task.Status = models.StatusError
		} else {
			task.Status = models.StatusComplete
			task.Content = content
		}

		if err := s.db.UpdateTask(task); err != nil {
			log.Printf("Error updating task %s: %v", taskID, err)
		}
	}
}
