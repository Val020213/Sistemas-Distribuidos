package main

import (
	"github.com/gin-gonic/gin"
)

type RequestStatus struct {
	httpStatus int
	body       gin.H
}

type WebPage struct {
	URL    string `json:"url" bson:"url"`
	Status string `json:"status" bson:"status"`
}
