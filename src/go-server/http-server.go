package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func buildApiResponse(statusCode int, status string, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"statusCode": statusCode,
		"status":     status,
		"message":    message,
		"data":       data,
	}
}

func getHome(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
}

func fetchURL(c *gin.Context) {

	url := c.Query("url")

	log.Println("Fetching URL: ", url)

	request := AddNewURL(url)

	if request.httpStatus != http.StatusOK {
		errorMessage, _ := request.body["error"].(string)
		c.JSON(request.httpStatus, buildApiResponse(request.httpStatus, "error", errorMessage, nil))
		return
	}

	c.JSON(request.httpStatus, buildApiResponse(request.httpStatus, "success", "URL fetched successfully", request.body))
}

func getUrlStates(c *gin.Context) {

	log.Println("Fetching URL states")

	request := listUrlStates()

	if request.httpStatus != http.StatusOK {
		errorMessage, _ := request.body["error"].(string)
		c.JSON(request.httpStatus, buildApiResponse(request.httpStatus, "error", errorMessage, nil))
		return
	}

	c.JSON(request.httpStatus, buildApiResponse(request.httpStatus, "success", "URL states fetched successfully", request.body))

}

func downloadURL(c *gin.Context) {

	url := c.Query("url")

	log.Println("Downloading URL: ", url)

	content, err := GetHTML(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/html", content.Bytes())
}

func main() {

	Init()
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/fetch", fetchURL)
	router.GET("/list", getUrlStates)
	router.GET("/download", downloadURL)
	router.GET("/", getHome)

	router.Run("localhost:8080")
}
