package api

import (
	"log"
	"scrapper_server/internal/models"
	"time"
)

func (s *Server) StartWorkerPool(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for taskID := range s.queue {
				log.Printf("Worker %d processing task: %s", id, taskID)

				// Actualizar estado a "in_progress"
				if err := s.db.UpdateTaskStatus(taskID, models.StatusInProgress); err != nil {
					log.Printf("Failed to update task status: %v", err)
					continue
				}

				// Simular procesamiento de scraping
				time.Sleep(5 * time.Second) // Aquí iría la lógica de scraping

				// Actualizar estado a "completed"
				if err := s.db.UpdateTaskStatus(taskID, models.StatusComplete); err != nil {
					log.Printf("Failed to update task status: %v", err)
				}

				log.Printf("Worker %d completed task: %s", id, taskID)
			}
		}(i)
	}
}
