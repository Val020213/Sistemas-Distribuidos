package api

import (
	"net/http"

	_ "scrapper_server/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Scrapper Server API
// @version 1.0
// @description This is a sample server for scrapping URLs.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// RegisterRoutes sets up the HTTP routes for the server using the Gin framework.
// It configures CORS settings to allow all origins, specific methods, and headers,
// and enables credentials for cookies/authentication.
//
// @Summary Register HTTP routes
// @Description Sets up the HTTP routes for the server including health check and task routes
// @Tags routes
// @Produce json
// @Success 200 {object} gin.Engine
// @Router / [get]
// @Router /health [get]
// @Router /tasks [post]
// @Router /tasks [get]

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cors.DefaultConfig().AllowOrigins, // Allow all origins for now
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// health check

	r.GET("/", s.HelloWorldHandler)

	// @Summary Health check
	// @Description Check the health of the server
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]

	r.GET("/health", s.healthHandler)

	// @Summary Create a task
	// @Description Create a new task to scrape a URL
	// @Tags tasks
	// @Accept json
	// @Produce json
	// @Param url body string true "URL to scrape"
	// @Success 201 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /tasks [post]

	r.POST("/tasks", s.CreateTaskHandler)

	// @Summary List all tasks
	// @Description Get a list of all tasks
	// @Tags tasks
	// @Produce json
	// @Success 200 {object} map[string][]models.Task
	// @Failure 500 {object} map[string]string
	// @Router /tasks [get]

	r.GET("/tasks", s.ListTasksHandler)

	// @Summary Get task content
	// @Description Get the content of a task by ID
	// @Tags tasks
	// @Produce json
	// @Param id path string true "Task ID"
	// @Success 200 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Router /tasks/{id}/content [get]

	r.GET("/scraper/tasks/:id/content", s.GetTaskContentHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
