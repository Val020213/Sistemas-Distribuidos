package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "Scrapper"
var bucketName = "HTMLs"

func ConnectDB() (*mongo.Client, *gridfs.Bucket) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create bucket
	bucket, err := gridfs.NewBucket(
		client.Database(bucketName),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	return client, bucket
}

func InitClient() (*mongo.Client, *gridfs.Bucket) {
	return ConnectDB()
}

func AddNewURL(client *mongo.Client, bucket *gridfs.Bucket, url string) RequestStatus {
	collection := client.Database(dbName).Collection("URLs")

	webPage := WebPage{
		URL:    url,
		Status: "pending",
	}

	_, err := collection.InsertOne(context.TODO(), webPage)
	if err != nil {
		return RequestStatus{
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		}
	}

	return RequestStatus{
		http.StatusOK,
		gin.H{"status": "URL added successfully"},
	}
}

func GetHTML(client *mongo.Client, bucket *gridfs.Bucket, url string) (bytes.Buffer, error) {

	var htmlContent bytes.Buffer

	downloadStream, err := bucket.OpenDownloadStreamByName(fmt.Sprintf("%s.html", url))
	if err != nil {
		return htmlContent, err
	}
	defer downloadStream.Close()

	_, err = io.Copy(&htmlContent, downloadStream)
	if err != nil {
		return htmlContent, err
	}

	return htmlContent, nil
}

func WorkerCall(client *mongo.Client, bucket *gridfs.Bucket, url string) {
	var worker Worker

	htmlContent, err := worker.Request(url)
	if err != nil {
		log.Fatal(err)
	}

	uploadStream, err := bucket.OpenUploadStream(fmt.Sprintf("%s.html", url))
	if err != nil {
		log.Fatal(err)
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(htmlContent.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// Change Status
	filter := bson.M{"url": url}
	update := bson.M{"$set": bson.M{"status": "scrapped"}}

	result, err := client.Database(dbName).Collection("URLs").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount == 0 {
		log.Fatal("No documents matched the filter")
	}
}
