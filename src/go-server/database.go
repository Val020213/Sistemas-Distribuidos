package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type URLStatus string

const (
	PendingStatus  URLStatus = "pending"
	ErrorStatus    URLStatus = "error"
	ScrappedStatus URLStatus = "scrapped"
)

var Scrapped map[string]bytes.Buffer //Add mutex
var Status map[string]URLStatus      //Add mutex
var myWorker *Worker

func Init() {
	Scrapped = make(map[string]bytes.Buffer)
	Status = make(map[string]URLStatus)
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

	Status[url] = PendingStatus

	go WorkerCall(url)

	return RequestStatus{
		http.StatusOK,
		gin.H{"status": "URL added successfully"},
	}
}

func listUrlStates() RequestStatus {
	var result []map[string]string

	i := 0
	for url, status := range Status {
		result = append(result, map[string]string{
			"id":     fmt.Sprintf("%d", i),
			"url":    url,
			"status": string(status),
		})
		i++
	}

	return RequestStatus{
		http.StatusOK,
		gin.H{"body": result},
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
		Status[url] = ErrorStatus
		return
	}

	Scrapped[url] = htmlContent
	Status[url] = ScrappedStatus
}
