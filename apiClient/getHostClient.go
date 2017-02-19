package apiClient

import (
	"net/url"

	"github.com/valyala/fasthttp"
)

func (client *Instance) getHostClient(addr *url.URL) (*fasthttp.HostClient, error) {
	client.clientsMu.RLock()
	if hostClient, ok := client.clients[addr.Host]; ok {
		client.clientsMu.RUnlock()
		return hostClient, nil
	}
	client.clientsMu.RUnlock()
	hostClient := &fasthttp.HostClient{
		Addr:  addr.Host,
		IsTLS: addr.Scheme == httpsScheme,
	}

	client.clientsMu.Lock()
	client.clients[addr.Host] = hostClient
	client.clientsMu.Unlock()
	return hostClient, nil
}
