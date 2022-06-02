package deploygate

import (
	go_deploygate "github.com/fnaoto/go_deploygate"
)

// Config : configuration for deploygate client
type Config struct {
	apiKey string
}

// Client : client for deploygate
type Client struct {
	client *go_deploygate.Client
}

// Client : API Client for deploygate
func (c *Config) Client() (interface{}, error) {
	var clnt Client

	client, err := go_deploygate.NewClient(c.apiKey)
	clnt.client = client
	return &clnt, err
}
