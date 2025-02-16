// file: mocks/task.go
package mocks

import (
	"server/internal/models"
	"time"
)

var Tasks = []models.TaskType{
	{

		URL:       "https://www.google.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{

		URL:       "https://www.facebook.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		URL:       "https://www.twitter.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{

		URL:       "https://www.linkedin.com",
		Status:    models.StatusInProgress,
		Content:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
