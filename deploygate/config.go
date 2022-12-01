package deploygate

import (
	go_deploygate "github.com/fnaoto/go_deploygate"
)

// Config : configuration for deploygate client
type Config struct {
	clientConfig go_deploygate.ClientConfig
	client       *go_deploygate.Client
}

// Client : API Client for deploygate
func (cfg *Config) initClient() error {
	c, err := go_deploygate.NewClient(cfg.clientConfig)
	cfg.client = c
	if err != nil {
		return err
	}
	return nil
}

func converEnterpriseMemberToMember(members []go_deploygate.EnterpriseMember) []go_deploygate.Member {
	var users []go_deploygate.Member

	for _, member := range members {
		user := go_deploygate.Member{
			Type:    member.Type,
			Name:    member.Name,
			IconUrl: member.IconUrl,
			Url:     member.Url,
		}
		users = append(users, user)
	}

	return users
}
