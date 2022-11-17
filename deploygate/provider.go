package deploygate

import (
	go_deploygate "github.com/fnaoto/go_deploygate"
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
			"deploygate_app_member":          dataSourceAppMember(),
			"deploygate_organization_member": dataSourceOrganizationMember(),
			"deploygate_enterprise_member":   dataSourceEnterpriseMember(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"deploygate_app_member":          resourceAppMember(),
			"deploygate_organization_member": resourceOrganizationMember(),
			"deploygate_enterprise_member":   resourceEnterpriseMember(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			clientConfig: go_deploygate.ClientConfig{
				ApiKey: d.Get("api_key").(string),
			},
		}

		err := config.initClient()
		if err != nil {
			return nil, err
		}

		return config.client, nil
	}
}
