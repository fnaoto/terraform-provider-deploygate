package deploygate

import (
	go_deploygate "github.com/recruit-mp/go-deploygate"
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
	client, err := go_deploygate.NewClient(c.apiKey)
	return client, err
}
