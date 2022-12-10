package deploygate

import (
	"context"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider : terraform provider
func Provider() func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("DG_API_KEY", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"deploygate_enterprise_member":              dataSourceEnterpriseMember(),
				"deploygate_enterprise_organization_member": dataSourceEnterpriseOrganizationMember(),
				"deploygate_organization_member":            dataSourceOrganizationMember(),
				"deploygate_organization_team_member":       dataSourceOrganizationTeamMember(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"deploygate_enterprise_member":              resourceEnterpriseMember(),
				"deploygate_enterprise_organization_member": resourceEnterpriseOrganizationMember(),
				"deploygate_organization_member":            resourceOrganizationMember(),
				"deploygate_organization_team_member":       resourceOrganizationTeamMember(),
			},
		}

		p.ConfigureContextFunc = providerConfigure(p)
		return p
	}
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config, err := initConfig(d)
		return config, diag.FromErr(err)
	}
}

func initConfig(d *schema.ResourceData) (*Config, error) {
	config := &Config{
		clientConfig: go_deploygate.ClientConfig{
			ApiKey: d.Get("api_key").(string),
		},
	}

	err := config.initClient()
	if err != nil {
		return nil, err
	}

	return config, nil
}
