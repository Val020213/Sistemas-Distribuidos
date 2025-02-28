package server

import (
	"fmt"
	"net/http"
	"server/internal/models"
	"time"

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

	r.POST("/task", s.getTaskByIDHandler)

	// <<<TASKS>>> END

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.node.Scraper.DB.Health())
}

func (s *Server) getTasksHandler(c *gin.Context) {
	tasks, err := s.node.Scraper.DB.GetTasks()
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
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "An error occurred while creating the task",
		})
		return
	}

	// Store task in chord ring

	task := models.TaskType{
		URL:       req.URL,
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.node.CallStoreData(task)

	// // Create task in DB
	// taskUrl, err := s.node.Scraper.DB.CreateTask(models.TaskType{
	// 	URL: req.URL,
	// })

	// if err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"statusCode": http.StatusInternalServerError,
	// 		"status":     "error",
	// 		"message":    "An error occurred while creating the task",
	// 	})
	// 	return
	// }
	// select {
	// case s.node.Scraper.TaskQueue <- taskUrl:
	// 	c.JSON(http.StatusAccepted, gin.H{
	// 		"statusCode": http.StatusOK,
	// 		"status":     "success",
	// 		"message":    "Task queued",
	// 		"data":       taskUrl,
	// 	})

	// default:
	// 	log.Printf("Queue full. Task ID: %v", taskUrl)
	// 	c.JSON(http.StatusServiceUnavailable, gin.H{
	// 		"statusCode": http.StatusServiceUnavailable,
	// 		"status":     "error",
	// 		"message":    "Task could not be queued",
	// 	})
	// }
}

func (s *Server) getTaskByIDHandler(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "An error occurred while getting the task",
		})
		return
	}

	task, err := s.node.Scraper.DB.GetTask(req.URL)

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
