package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Scrapped map[string]bytes.Buffer //Add mutex
var Status map[string]string         //Add mutex
var myWorker *Worker

func Init() {
	Scrapped = make(map[string]bytes.Buffer)
	Status = make(map[string]string)
	myWorker = &Worker{}
}

func AddNewURL(url string) RequestStatus {

	_, exists := Status[url]

	if exists {
		return RequestStatus{
			http.StatusBadRequest,
			gin.H{"error": "URL already in system"},
		}
	}

	Status[url] = "pending"

	go WorkerCall(url)

	return RequestStatus{
		http.StatusOK,
		gin.H{"status": "URL added successfully"},
	}
}

func GetStates() RequestStatus {

	return RequestStatus{
		http.StatusOK,
		gin.H{"body": Status},
	}

}

type HTMLError string

func (e HTMLError) Error() string {
	return fmt.Sprintf("The following url %s is not ready", string(e))
}

func GetHTML(url string) (bytes.Buffer, error) {

	htmlContent, exists := Scrapped[url]

	if !exists {
		return htmlContent, HTMLError(url)
	}

	return htmlContent, nil
}

func WorkerCall(url string) {

	htmlContent, err := myWorker.Request(url)
	if err != nil {
		Scrapped[url] = htmlContent
		Status[url] = "error"
		return
	}

	Scrapped[url] = htmlContent
	Status[url] = "finish"
}
