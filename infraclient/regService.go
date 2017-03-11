package infraclient

import (
	"github.com/gramework/gramework/apiClient"
	"github.com/gramework/gramework/infrastructure"
)

type InfraAPI struct {
	URL    string
	client *apiClient.Instance
}

func New(url string) *InfraAPI {
	return &InfraAPI{
		URL: url,
		client: apiClient.New(apiClient.Config{
			Addresses: []string{url},
		}),
	}
}

func (i *InfraAPI) RegisterService(s infrastructure.Service) error { return nil }
