package main

import (
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Worker struct {
	mu sync.Mutex
}

func (w *Worker) Request(url string) RequestStatus {
	w.mu.Lock()
	defer w.mu.Unlock()

	log.Println("Fetching URL: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return RequestStatus{
			http.StatusBadRequest,
			gin.H{
				"error":  err.Error(),
				"detail": "Failed to fetch the URL",
			},
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RequestStatus{
			http.StatusConflict,
			gin.H{
				"error":  err.Error(),
				"detail": "Failed to read the response body",
			},
		}
	}

	return RequestStatus{
		http.StatusOK,
		gin.H{
			"body": string(body),
		},
	}
}
