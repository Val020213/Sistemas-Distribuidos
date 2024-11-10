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

func (w *Worker) Request(url string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	log.Println("Fetching URL: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
