package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"
)

type Worker struct {
	mu sync.Mutex
}

func (w *Worker) Request(url string) (bytes.Buffer, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var htmlContent bytes.Buffer

	log.Println("Fetching URL: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return htmlContent, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(&htmlContent, resp.Body)
	if err != nil {
		return htmlContent, err
	}

	return htmlContent, nil
}
