// scraper/scraper.go
package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"
)

type Scraper struct {
	DB        database.Service
	TaskQueue chan string
}

func NewScraper() *Scraper {
	queueSize := utils.GetEnvAsInt("TASK_QUEUE_SIZE", 100)
	numWorkers := utils.GetEnvAsInt("NUM_WORKERS", 5)

	scraper := &Scraper{
		DB:        database.New(),
		TaskQueue: make(chan string, queueSize),
	}

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		go scraper.worker()
	}

	return scraper
}

func (s *Scraper) worker() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker recovered from panic: %v", r)
		}
	}()

	for taskID := range s.TaskQueue {

		taskIDUint, err := strconv.ParseUint(taskID, 10, 32)
		if err != nil {
			log.Printf("Invalid task ID format %s: %v", taskID, err)
			continue
		}

		task, err := s.DB.GetTask(uint32(taskIDUint))
		if err != nil {
			log.Printf("Error fetching task %d: %v", taskIDUint, err)
			continue
		}

		// Scrape URL
		content, err := scrapeURL(task.URL)
		if err != nil {
			log.Printf("Error scraping URL %s: %v", task.URL, err)
			task.Status = models.StatusError
		} else {
			task.Status = models.StatusComplete
			task.Content = content
		}

		if err := s.DB.UpdateTask(task); err != nil {
			log.Printf("Error updating task %s: %v", taskID, err)
		}
	}
}

func scrapeURL(url string) (string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to scrape URL: %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
