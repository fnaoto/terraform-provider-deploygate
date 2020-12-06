package deploygate

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	go_deploygate "github.com/recruit-mp/go-deploygate"
)

func dataSourceAppCollaborator() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppCollaboratorRead,

		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAppCollaboratorRead(d *schema.ResourceData, meta interface{}) error {
	client := go_deploygate.DefaultClient()

	owner := d.Get("owner").(string)
	platform := d.Get("platform").(string)
	appID := d.Get("app_id").(string)

	g := &go_deploygate.GetAppCollaboratorInput{
		Owner:    owner,
		Platform: platform,
		AppId:    appID,
	}
	collaborator, _ := client.GetAppCollaborator(g)
	d.SetId(fmt.Sprintf("%s/%s/%s", owner, platform, appID))
	d.Set("users", collaborator.Results.Users)
	d.Set("teams", collaborator.Results.Teams)
	d.Set("usage", collaborator.Results.Usage)
	return nil
}
