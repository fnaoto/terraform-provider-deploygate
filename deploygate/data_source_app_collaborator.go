package deploygate

import (
	"fmt"
	"log"

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
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"teams": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"usage": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     schema.TypeInt,
			},
		},
	}
}

func dataSourceAppCollaboratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).client

	owner := d.Get("owner").(string)
	platform := d.Get("platform").(string)
	appID := d.Get("app_id").(string)

	log.Printf("[DEBUG] dataSourceAppCollaboratorRead: %s, %s, %s", owner, platform, appID)

	g := &go_deploygate.GetAppCollaboratorInput{
		Owner:    owner,
		Platform: platform,
		AppId:    appID,
	}

	collaborator, err := client.GetAppCollaborator(g)

	if err != nil {
		return err
	}

	rs := collaborator.Results

	d.SetId(fmt.Sprintf("%s/%s/%s", owner, platform, appID))
	d.Set("users", rs.Users)
	d.Set("teams", rs.Teams)
	d.Set("usage", map[string]interface{}{
		"max":  rs.Usage.Max,
		"used": rs.Usage.Used,
	})

	return nil
}
