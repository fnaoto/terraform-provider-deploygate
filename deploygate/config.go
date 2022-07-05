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
func (cfg *Config) Client() (*Client, error) {
	var clnt Client

	c, err := go_deploygate.NewClient(cfg.apiKey)
	clnt.client = c

	return &clnt, err
}
