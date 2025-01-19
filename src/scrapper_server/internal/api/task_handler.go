package api

import (
	"net/http"
	"scrapper_server/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) CreateTaskHandler(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &models.Task{
		ID:        uuid.NewString(),
		URL:       req.URL,
		Status:    models.StatusInProgress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	s.workerPool.queue <- task.ID

	s.workerPool.RestartWorkerPool()

	c.JSON(http.StatusCreated, gin.H{"taskId": task.ID, "status": task.Status})
}

func (s *Server) ListTasksHandler(c *gin.Context) {
	// Obtener todas las tareas del repositorio
	tasks, err := s.db.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Devolver las tareas en formato JSON
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
