package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"server/internal/models"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	CreateTask(task models.TaskType) (string, error)
	UpdateTask(task models.TaskType) error
	GetTasks() ([]models.TaskType, error)
	GetTask(id string) (models.TaskType, error)
}

type service struct {
	db *mongo.Client
}

var (
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	database = os.Getenv("BLUEPRINT_DB_NAME")
	username = os.Getenv("BLUEPRINT_DB_USER")
	password = os.Getenv("BLUEPRINT_DB_PASSWORD")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)))

	if err != nil {
		log.Fatal(err)

	}

	// Create Collections here
	client.Database(database).CreateCollection(context.Background(), "tasks")

	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down agua: %v", err)
	}
	envVars := os.Environ()
	for _, env := range envVars {
		fmt.Println(env)
	}
	return map[string]string{
		"message": "It's healthy",
	}
}

// Task Repository

func (s *service) CreateTask(task models.TaskType) (string, error) {
	// Set a context with timeout for the insert operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set the timestamps for task creation and last update.
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	task.Status = models.StatusInProgress
	if task.ID == "" {
		task.ID = primitive.NewObjectID().Hex()
	}

	// Get the "tasks" collection from the configured database.
	collection := s.db.Database(database).Collection("tasks")

	// Insert the task into the collection.
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to insert task: %v", err)
	}

	return task.ID, nil
}

func (s *service) UpdateTask(task models.TaskType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")

	filter := bson.M{"_id": task.ID}
	update := bson.M{"$set": bson.M{
		"status":     task.Status,
		"content":    task.Content,
		"updated_at": time.Now(),
	}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *service) GetTasks() ([]models.TaskType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.TaskType
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *service) GetTask(id string) (models.TaskType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")

	filter := bson.M{"_id": id}

	var task models.TaskType
	err := collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return models.TaskType{}, err
	}

	return task, nil
}
