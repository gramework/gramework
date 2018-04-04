package infrastructure

import (
	"sync"
	"time"
)

// New initializes an empty infrastructure
func New() *Infrastructure {
	return &Infrastructure{
		Lock:            new(sync.RWMutex),
		Services:        make(map[string]*Service),
		UpdateTimestamp: time.Now().UnixNano(),
	}
}
