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
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DG_USER_NAME", nil),
			},
			"organization_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DG_ORGANIZATION_NAME", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"deploygate_app_collaborator": dataSourceAppCollaborator(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"deploygate_app_collaborator": resourceAppCollaborator(),
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
