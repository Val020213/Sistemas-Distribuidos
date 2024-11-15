package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

var dbClient *mongo.Client
var dbBucket *gridfs.Bucket

func getHome(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nil)
}

// fetchURL fetches the content of the URL passed as a query parameter
// example : http://127.0.0.1:8080/fetch?url=https://example.com
func fetchURL(c *gin.Context) {

	url := c.Query("url")

	log.Println("Fetching URL: ", url)

	request := AddNewURL(dbClient, dbBucket, url)

	c.JSON(request.httpStatus, request.body)
}

func downloadURL(c *gin.Context) {

	url := c.Query("url")

	log.Println("Downloading URL: ", url)

	content, err := GetHTML(dbClient, dbBucket, url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/html", content.Bytes())
}

func main() {

	dbClient, dbBucket = InitClient()

	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
		log.Println("Connection to MongoDB closed.")
	}()

	router := gin.Default()

	router.GET("/fetch", fetchURL)
	router.GET("/download", downloadURL)
	router.GET("/", getHome)

	router.Run("127.0.0.1:8080")

}
