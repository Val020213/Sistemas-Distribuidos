// scraper/scraper.go
package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"server/internal/database"
	"server/internal/models"
	"server/internal/utils"
)

type Scraper struct {
	DB          database.Service
	TaskQueue   chan uint64
	TaskWorking map[uint64]bool
}

func NewScraper() *Scraper {
	queueSize := utils.GetEnvAsInt("TASK_QUEUE_SIZE", 100)
	numWorkers := utils.GetEnvAsInt("NUM_WORKERS", 5)

	scraper := &Scraper{
		DB:          database.New(),
		TaskQueue:   make(chan uint64, queueSize),
		TaskWorking: make(map[uint64]bool),
	}

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		go scraper.worker()

	}

	go func() {
		interval := utils.GetEnvAsInt("TASK_QUEUE_INTERVAL_SECONDS", 5)
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for range ticker.C {
			scraper.populateTaskQueue()
		}
	}()

	return scraper
}

func (s *Scraper) populateTaskQueue() {

	if len(s.TaskQueue) > 0 || len(s.TaskWorking) > 0 {
		return
	}

	pendingTasks, err := s.DB.GetTasksWithFilter(map[string]interface{}{"status": models.StatusInProgress})

	if err != nil {
		utils.RedPrint("Error fetching pending tasks: ", err)
		return
	}

	var wg sync.WaitGroup
	for _, task := range pendingTasks {
		wg.Add(1)
		go func(taskKey uint64) {
			defer wg.Done()
			s.TaskQueue <- taskKey
		}(task.Key)
	}
	wg.Wait()

	log.Println("Task queue populated with pending tasks.")
}

func (s *Scraper) worker() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Worker recovered from panic: %v", r)
		}
	}()

	for taskID := range s.TaskQueue {
		s.TaskWorking[taskID] = true
		task, err := s.DB.GetTask(taskID)
		if err != nil {
			delete(s.TaskWorking, taskID)
			log.Printf("Error fetching task %v: %v", taskID, err)
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
			task.UpdatedAt = time.Now()
		}

		if err := s.DB.UpdateTask(task); err != nil {
			log.Printf("Error updating task with id %v: %v", taskID, err)
		}

		delete(s.TaskWorking, taskID)
	}
}

func scrapeURL(url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scrape URL: %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
