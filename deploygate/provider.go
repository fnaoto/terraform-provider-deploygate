package deploygate

import (
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
			"deploygate_app_collaborator":    dataSourceAppCollaborator(),
			"deploygate_organization_member": dataSourceOrganizationMember(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"deploygate_app_collaborator":    resourceAppCollaborator(),
			"deploygate_organization_member": resourceOrganizationMember(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			apiKey: d.Get("api_key").(string),
		}

		meta, err := config.Client()

		if err != nil {
			return nil, err
		}

		return meta, nil
	}
}
