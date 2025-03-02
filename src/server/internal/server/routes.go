package server

import (
	"fmt"
	"net/http"
	"server/internal/models"
	"server/internal/utils"
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
	utils.GreenPrint("GET Tasks Handler123")
	tasks, err := s.node.CallList()
	utils.GreenPrint("CALL LIST: ", tasks)
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
	utils.GreenPrint("Create Task Handler")
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
		Key:       uint64(utils.ChordHash(req.URL, s.node.M)),
		URL:       req.URL,
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.node.CallCreateData(task)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occurred while creating the task",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "Task queued",
		"data":       task,
	})
}

func (s *Server) getTaskByIDHandler(c *gin.Context) {
	utils.GreenPrint("Get Task by ID Handler")

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

	data, err := s.node.CallGetData(req.URL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "Task fetched successfully",
		"data":       data,
	})

}
