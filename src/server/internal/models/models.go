package models

import (
	"errors"
	"net/url"
	"time"
)

type TaskStatusType string

const (
	StatusInProgress TaskStatusType = "in_progress"
	StatusComplete   TaskStatusType = "complete"
	StatusError      TaskStatusType = "error"
)

type TaskType struct {
	ID        string         `bson:"_id" json:"id"`
	URL       string         `bson:"url" json:"url"`
	Status    TaskStatusType `bson:"status" json:"status"`
	Content   string         `bson:"content,omitempty" json:"content,omitempty"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at" json:"updated_at"`
}

func (t *TaskType) Validate() error {
	if t.URL == "" {
		return errors.New("URL cannot be empty")
	}
	if _, err := url.ParseRequestURI(t.URL); err != nil {
		return errors.New("invalid URL format")
	}
	if t.Status == "" {
		return errors.New("status cannot be empty")
	}
	return nil
}
