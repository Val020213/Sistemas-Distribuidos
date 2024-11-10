package main

import (
	"github.com/gin-gonic/gin"
)

type RequestStatus struct {
	httpStatus int
	body       gin.H
}
