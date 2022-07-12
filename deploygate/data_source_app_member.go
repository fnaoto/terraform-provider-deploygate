package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppMemberRead,
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
							Default:  2,
						},
					},
				},
			},
		},
	}
}

func dataSourceAppMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).client

	owner := d.Get("owner").(string)
	platform := d.Get("platform").(string)
	appID := d.Get("app_id").(string)

	log.Printf("[DEBUG] dataSourceAppMemberRead: %s, %s, %s", owner, platform, appID)

	g := &go_deploygate.GetAppMembersRequest{
		Owner:    owner,
		Platform: platform,
		AppId:    appID,
	}

	member, err := client.GetAppMembers(g)

	if err != nil {
		return err
	}

	rs := member.Results

	d.SetId(fmt.Sprintf("%s/%s/%s", owner, platform, appID))
	d.Set("users", rs.Users)

	return nil
}
