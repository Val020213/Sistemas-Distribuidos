package server

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/models"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Add your frontend URL
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"Accept", "Authorization", "Content-Type"},
		// AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	// <<<TASKS>>> START

	r.GET("/tasks", s.getTasksHandler)

	r.POST("/tasks", s.createTaskHandler)

	r.GET("/task/:id", s.getTaskByIDHandler)

	// <<<TASKS>>> END

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) getTasksHandler(c *gin.Context) {
	tasks, err := s.db.GetTasks()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occurred while fetching tasks",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "Tasks fetched successfully",
		"data":       tasks,
	})
}

func (s *Server) createTaskHandler(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occurred while creating the task",
		})
		return
	}

	// Create task in DB
	taskID, err := s.db.CreateTask(models.TaskType{
		URL: req.URL,
	})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occurred while creating the task",
		})
		return
	}
	select {
	case s.TaskQueue <- strconv.FormatUint(uint64(taskID), 10):
		c.JSON(http.StatusAccepted, gin.H{
			"statusCode": http.StatusOK,
			"status":     "success",
			"message":    "Task queued",
			"data":       taskID,
		})

	default:
		log.Printf("Queue full. Task ID: %v", taskID)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"statusCode": http.StatusServiceUnavailable,
			"status":     "error",
			"message":    "Task could not be queued",
		})
	}
}

func (s *Server) getTaskByIDHandler(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id path parameter is required"})
		return
	}

	taskIDUint, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}
	task, err := s.db.GetTask(uint32(taskIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "Task fetched successfully",
		"data":       task,
	})
}
