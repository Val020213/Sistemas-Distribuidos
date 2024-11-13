package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbName = "Scrapper"

func ConnectDB() *mongo.Client {
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

	log.Println("Connected to MongoDB!")

	return client
}

func InitClient() *mongo.Client {
	return ConnectDB()
}

func AddNewURL(client *mongo.Client, url string) RequestStatus {
	collection := client.Database(dbName).Collection("URLs")

	webPage := WebPage{
		URL:     url,
		Status:  "pending",
		Content: "not available",
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


