package deploygate

import (
	go_deploygate "github.com/fnaoto/go_deploygate"
)

// Config : configuration for deploygate client
type Config struct {
	apiKey string
	client *go_deploygate.Client
}

// Client : API Client for deploygate
func (cfg *Config) initClient() error {
	c, err := go_deploygate.NewClient(cfg.apiKey)
	cfg.client = c
	if err != nil {
		return err
	}
	return nil
}
