package infrastructure

import (
	"sync"
)

type (
	// Infrastructure handles lists of services and current timestamp
	Infrastructure struct {
		Services         map[string]*Service `json:"services"`
		UpdateTimestamp  int64               `json:"update_timestamp"`
		CurrentTimestamp int64               `json:"current_timestamp"`
		Lock             *sync.RWMutex       `json:"-"`
	}

	// Service defines an abstract service in infrastructure
	Service struct {
		Name      string      `json:"name"`
		Addresses []Address   `json:"addresses"`
		Type      ServiceType `json:"type"`
	}

	// Address defines service addr
	Address struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		URL  string `json:"url"`
	}

	// ServiceType defines a type of protocol that
	// service using
	ServiceType string
)
