package models

import (
	"net/http"
	"time"
)

type DiscoveredServer struct {
	Address  string
	LastSeen time.Time
}

type ServerResponse struct {
	Resp *http.Response
	Err  error
}
