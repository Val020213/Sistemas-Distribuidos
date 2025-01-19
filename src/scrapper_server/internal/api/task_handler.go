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

	s.queue <- task.ID

	c.JSON(http.StatusCreated, gin.H{"taskId": task.ID, "status": task.Status})
}
