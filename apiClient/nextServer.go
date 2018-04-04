package apiClient

import (
	"net/url"

	"github.com/gramework/gramework"
)

func (client *Instance) nextServer() (*requestInfo, error) {
	if len(client.conf.Addresses) == 0 {
		return nil, ErrNoServerAvailable
	}

	for i := 0; i < len(client.conf.Addresses); i++ {
		addr := client.conf.Addresses[client.balancer.next()]
		hostURL, err := url.Parse(addr)
		if err != nil {
			gramework.Errorf("error while parsing host url: %s", err)
			continue
		}

		return &requestInfo{
			HostClient: client.getHostClient(hostURL),
			Addr:       addr,
		}, nil
	}

	return nil, ErrNoServerAvailable
}
