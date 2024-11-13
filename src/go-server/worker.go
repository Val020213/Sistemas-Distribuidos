package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type Worker struct {
	mu sync.Mutex
}

func (w *Worker) Request(url string) WebPage {
	w.mu.Lock()
	defer w.mu.Unlock()

	log.Println("Fetching URL: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return WebPage{
			URL:     url,
			Status:  "failed",
			Content: err.Error(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WebPage{
			URL:     url,
			Status:  "failed",
			Content: err.Error(),
		}
	}

	return WebPage{
		URL:     url,
		Status:  "success",
		Content: string(body),
	}

}
