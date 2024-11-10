package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHome(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
}

// fetchHTTP fetches the content of the URL passed as a query parameter
// example : http://127.0.0.1:8080/fetch?url=https://example.com
func fetchHTTP(c *gin.Context) {

	url := c.Query("url")

	log.Println("Fetching URL: ", url)

	var w Worker
	body, err := w.Request(url)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"detail": "Failed to fetch the URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"body": body,
	})
}

func main() {

	router := gin.Default()

	router.GET("/fetch", fetchHTTP)
	router.GET("/", getHome)

	router.Run("127.0.0.1:8080")

}
