package infrastructure

import (
	"errors"
)

var (
	// ErrServiceTypeNotFound occurs when unknown service type
	// passed in GetTypeByString
	ErrServiceTypeNotFound = errors.New("service type not found")
)

// GetTypeByString returns a ServiceType associated with
// given typeName and nil, if type is known, otherwise
// instantiates a ServiceType with given typeName
// but returns ErrServiceTypeNotFound
func GetTypeByString(typeName string) (ServiceType, error) {
	switch ServiceType(typeName) {
	case HTTP:
		return HTTP, nil
	case HTTPS:
		return HTTPS, nil
	case TCP:
		return TCP, nil
	case UDP:
		return UDP, nil
	case Custom:
		return Custom, nil
	default:
		return ServiceType(typeName), ErrServiceTypeNotFound
	}
}
