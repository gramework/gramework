package apiClient

import (
	"bytes"
	"encoding/json"
)

// GET sends a request with GET method
func (client *Instance) GET() (statusCode int, body []byte, err error) {
	api, err := client.nextServer()
	if err != nil {
		return 0, nil, err
	}
	return api.HostClient.Get(nil, api.Addr)
}

// GetJSON sends a GET request and deserializes response in a provided variable
func (client *Instance) GetJSON(v interface{}) (statusCode int, err error) {
	statusCode, body, err := client.GET()
	if err != nil {
		return 0, err
	}

	return statusCode, json.NewDecoder(bytes.NewReader(body)).Decode(&v)
}
