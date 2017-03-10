package infrastructure

import "errors"

var (
	// ErrServiceExists occurs when you trying to register a service that already exists
	ErrServiceExists = errors.New("service exists")
)

// RegisterService in the infrastructure
func (i *Infrastructure) RegisterService(name string, s Service) error {
	i.Lock.RLock()
	if _, ok := i.Services[name]; ok {
		i.Lock.RUnlock()
		return ErrServiceExists
	}
	i.Lock.RUnlock()
	i.Lock.Lock()
	if s.Name == "" {
		s.Name = name
	}
	i.Services[name] = &s
	i.Lock.Unlock()

	return nil
}

// RegisterServiceBatch registers services from given map[service name]Service
// Returns error if service already exists
func (i *Infrastructure) RegisterServiceBatch(m map[string]Service) error {
	for name, s := range m {
		err := i.RegisterService(name, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterServiceBatchIgnore registers services from given map[service name]Service
// Ignores any error if service already exists
func (i *Infrastructure) RegisterServiceBatchIgnore(m map[string]Service) {
	for name, s := range m {
		i.RegisterService(name, s)
	}
}

// MergeService in the infrastructure
func (i *Infrastructure) MergeService(name string, s Service) {
	i.Lock.Lock()
	if _, ok := i.Services[name]; ok {
		i.Services[name].Addresses = append(i.Services[name].Addresses, s.Addresses...)
		i.Lock.Unlock()
		return
	}
	i.Services[name] = &s
	i.Lock.Unlock()
}
