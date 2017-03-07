package infrastructure

import "sync"

// Infrastructure handles lists of services and current timestamp
type Infrastructure struct {
	Services         map[string]*Service `json:"services"`
	UpdateTimestamp  int64               `json:"update_timestamp"`
	CurrentTimestamp int64               `json:"current_timestamp"`
	Lock             *sync.RWMutex       `json:"-"`
}

// Service defines an abstract service in infrastructure
type Service struct {
	Addresses []Address   `json:"addresses"`
	Type      ServiceType `json:"type"`
}

// Address defines service addr
type Address struct {
	Host string `json:"host"`
	Port int    `json:"post"`
}

// ServiceType defines a type of protocol that
// service using
type ServiceType string
