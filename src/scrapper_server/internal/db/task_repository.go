package database

import (
	"context"
	"time"

	"scrapper_server/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskService define las operaciones sobre las tareas
type TaskService interface {
	CreateTask(task *models.Task) error
	UpdateTaskStatus(id string, status models.TaskStatus) error
	GetTaskByID(id string) (*models.Task, error)
}

// taskService implementación de TaskService
type taskService struct {
	collection *mongo.Collection
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService(client *mongo.Client) TaskService {
	return &taskService{
		collection: client.Database("scrapper").Collection("tasks"),
	}
}

// CreateTask inserta una nueva tarea en la base de datos
func (s *taskService) CreateTask(task *models.Task) error {
	_, err := s.collection.InsertOne(context.Background(), task)
	return err
}

// UpdateTaskStatus actualiza el estado de una tarea existente
func (s *taskService) UpdateTaskStatus(id string, status models.TaskStatus) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}
	_, err := s.collection.UpdateOne(context.Background(), filter, update)
	return err
}

// GetTaskByID obtiene una tarea por su ID
func (s *taskService) GetTaskByID(id string) (*models.Task, error) {
	var task models.Task
	err := s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
