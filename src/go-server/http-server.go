package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbClient *mongo.Client

func getHome(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
}

// fetchHTTP fetches the content of the URL passed as a query parameter
// example : http://127.0.0.1:8080/fetch?url=https://example.com
func fetchHTTP(c *gin.Context) {

	url := c.Query("url")

	log.Println("Fetching URL: ", url)

	request := AddNewURL(dbClient, url)

	c.JSON(request.httpStatus, request.body)
}

func main() {

	dbClient = InitClient()

	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
		log.Println("Connection to MongoDB closed.")
	}()

	router := gin.Default()

	router.GET("/fetch", fetchHTTP)
	router.GET("/", getHome)

	router.Run("127.0.0.1:8080")

}
