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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	CreateTask(task models.TaskType) (string, error)
	UpdateTask(task models.TaskType) error
	GetTasksWithFilter(filter bson.M) ([]models.TaskType, error)
	GetTasks() ([]models.TaskType, error)
	GetTask(key uint64) (models.TaskType, error)
	DeleteData(filter bson.M) error
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
	// Set a context with timeout for the operation.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	task.Status = models.StatusInProgress

	collection := s.db.Database(database).Collection("tasks")

	filter := bson.M{
		"key": task.Key,
	}

	update := bson.M{
		"$set": bson.M{
			"content":    task.Content,
			"created_at": now,
			"updated_at": now,
			"status":     models.StatusInProgress,
		},
	}

	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", fmt.Errorf("failed to update existing task: %v", err)
	}

	fmt.Println("update result", updateResult)
	if updateResult.MatchedCount > 0 {
		return task.URL, nil
	}

	fmt.Println("Inserting...")
	_, err = collection.InsertOne(ctx, task)
	if err != nil {
		return "", fmt.Errorf("failed to insert task: %v", err)
	}

	return task.URL, nil
}

func (s *service) UpdateTask(task models.TaskType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")

	filter := bson.M{"key": task.Key}

	update := bson.M{"$set": bson.M{
		"url":        task.URL,
		"key":        task.Key,
		"status":     task.Status,
		"content":    task.Content,
		"created_at": task.CreatedAt,
		"updated_at": task.UpdatedAt,
	}}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *service) GetTasksWithFilter(filter bson.M) ([]models.TaskType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")
	cursor, err := collection.Find(ctx, filter)
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

func (s *service) GetTasks() ([]models.TaskType, error) {
	return s.GetTasksWithFilter(bson.M{})
}

func (s *service) GetTask(key uint64) (models.TaskType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")

	filter := bson.M{"key": key}

	var task models.TaskType
	err := collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return models.TaskType{}, err
	}

	return task, nil
}

func (s *service) DeleteData(filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.db.Database(database).Collection("tasks")

	_, err := collection.DeleteMany(ctx, filter)
	return err
}
