package main

import (
	"io"
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

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"detail": "Failed to fetch the URL",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"detail": "Failed to read the response body",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"body": string(body),
	})
}

func main() {

	router := gin.Default()

	router.GET("/fetch", fetchHTTP)
	router.GET("/", getHome)

	router.Run("127.0.0.1:8080")

}
