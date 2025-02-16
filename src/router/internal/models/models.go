package models

import "time"

type DiscoveredServer struct {
	Address  string
	LastSeen time.Time
}
