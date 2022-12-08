package deploygate

import (
	"context"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider : terraform provider
func Provider() *schema.Provider {
	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DG_API_KEY", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"deploygate_organization_team_member": dataSourceOrganizationTeamMember(),
			"deploygate_organization_member":      dataSourceOrganizationMember(),
			"deploygate_enterprise_member":        dataSourceEnterpriseMember(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"deploygate_organization_team_member": resourceOrganizationTeamMember(),
			"deploygate_organization_member":      resourceOrganizationMember(),
			"deploygate_enterprise_member":        resourceEnterpriseMember(),
		},
	}

	p.ConfigureContextFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := &Config{
			clientConfig: go_deploygate.ClientConfig{
				ApiKey: d.Get("api_key").(string),
			},
		}

		err := config.initClient()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return config, nil
	}
}
