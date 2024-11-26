package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHome(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
}

func fetchURL(c *gin.Context) {

	url := c.Query("url")

	log.Println("Fetching URL: ", url)

	request := AddNewURL(url)

	c.JSON(request.httpStatus, request.body)
}

func getStates(c *gin.Context) {

	log.Println("Fetching URL states")

	request := GetStates()

	c.JSON(request.httpStatus, request.body)
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

	router.GET("/fetch", fetchURL)
	router.GET("/state", getStates)
	router.GET("/download", downloadURL)
	router.GET("/", getHome)

	router.Run("localhost:8080")
}
