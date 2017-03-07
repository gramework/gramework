package infrastructure

import (
	"errors"
	"time"
)

var (
	// ErrServiceNotExists occurs when you trying to add a server to service that not exists
	ErrServiceNotExists = errors.New("service is not exists")
)

// AddServiceServer registers server to a service in the infrastructure
func (i *Infrastructure) AddServiceServer(serviceName string, addr Address) error {
	i.Lock.RLock()
	if _, ok := i.Services[serviceName]; !ok {
		i.Lock.RUnlock()
		return ErrServiceNotExists
	}
	i.Lock.RUnlock()
	i.Lock.Lock()
	if i.Services[serviceName].Addresses == nil {
		i.Services[serviceName].Addresses = make([]Address, 0)
	}
	i.UpdateTimestamp = time.Now().UnixNano()
	i.Services[serviceName].Addresses = append(i.Services[serviceName].Addresses, addr)
	i.Lock.Unlock()
	return nil
}
