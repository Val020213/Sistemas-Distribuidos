// file: mocks/task.go
package mocks

import (
	"server/internal/models"
	"time"
)

var Tasks = []models.TaskType{
	{
		ID:        "1",
		URL:       "https://www.google.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "2",
		URL:       "https://www.facebook.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "3",
		URL:       "https://www.twitter.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "4",
		URL:       "https://www.linkedin.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
