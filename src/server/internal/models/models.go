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
	URL       string         `bson:"url" json:"url"`
	Key       uint64         `bson:"key" json:"key"`
	Status    TaskStatusType `bson:"status" json:"status"`
	Content   []byte         `bson:"content,omitempty" json:"content,omitempty"`
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
