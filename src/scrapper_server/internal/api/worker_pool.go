package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"scrapper_server/internal/models"
	"sync"
	"time"
)

type workerPool struct {
	workerCount int
	queue       chan string
	mu          sync.Mutex
}

func NewWorkerPool(workerCount, workerMaxTask int) workerPool {
	return workerPool{
		workerCount: workerCount,
		mu:          sync.Mutex{},
		queue:       make(chan string, workerMaxTask),
	}
}

func (wp *workerPool) RestartWorkerPool() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	close(wp.queue)
	wp.queue = make(chan string, cap(wp.queue))
}

func (s *Server) StartWorkerPool(workerCount int) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout global para todas las solicitudes HTTP
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			s.worker(ctx, id, client)
		}(i)
	}

	// Esperar a que todos los workers terminen
	wg.Wait()
}

func (s *Server) worker(ctx context.Context, id int, client *http.Client) {
	for {
		select {
		case taskID, ok := <-s.workerPool.queue:
			if !ok {
				return
			}
			log.Printf("Worker %d processing task: %s", id, taskID)
			s.processTask(client, taskID)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Server) processTask(client *http.Client, taskID string) {
	task, err := s.db.GetTaskByID(taskID)
	if err != nil {
		log.Printf("Failed to get task: %v", err)
		s.updateTaskStatus(taskID, models.StatusError)
		return
	}

	content, err := GetUrlContent(client, task.URL)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", task.URL, err)
		s.updateTaskStatus(taskID, models.StatusError)
		return
	}

	err = s.db.SaveTaskContent(taskID, content)
	if err != nil {
		log.Printf("Failed to save content for task %s: %v", taskID, err)
		s.updateTaskStatus(taskID, models.StatusError)
		return
	}

	s.updateTaskStatus(taskID, models.StatusComplete)
	log.Printf("Worker completed task: %s", taskID)
}

func (s *Server) updateTaskStatus(taskID string, status models.TaskStatus) {
	if err := s.db.UpdateTaskStatus(taskID, status); err != nil {
		log.Printf("Failed to update task status to %v: %v", status, err)
	}
}

// GetUrlContent obtiene el contenido de una URL usando un cliente HTTP reutilizable
func GetUrlContent(client *http.Client, url string) (string, error) {
	// Crear un contexto con timeout específico para esta solicitud
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Realizar la solicitud HTTP
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Verificar que la respuesta sea exitosa
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Usar io.ReadAll para leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body) // Actualización aquí, usando io.ReadAll
	if err != nil {
		return "", err
	}

	return string(body), nil
}
