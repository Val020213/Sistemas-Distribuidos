package database

import (
	"context"
	"os"
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
	GetAllTasks() ([]*models.Task, error)
	SaveTaskContent(taskID string, content string) error
}

// taskService implementación de TaskService
type taskService struct {
	collection *mongo.Collection
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService(client *mongo.Client) TaskService {
	return &taskService{
		collection: client.Database(os.Getenv("BLUEPRINT_DB_NAME")).Collection("tasks"),
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

// GetAllTasks obtiene todas las tareas de la base de datos
func (s *taskService) GetAllTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	cursor, err := s.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SaveTaskContent guarda el contenido descargado en la base de datos para una tarea específica.
func (s *taskService) SaveTaskContent(taskID string, content string) error {
	filter := bson.M{"_id": taskID}
	update := bson.M{"$set": bson.M{"content": content, "updated_at": time.Now()}}
	_, err := s.collection.UpdateOne(context.Background(), filter, update)
	return err
}
