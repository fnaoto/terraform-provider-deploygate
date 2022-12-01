package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganizationTeamMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationTeamMemberRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationTeamMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).client

	organization := d.Get("organization").(string)
	team := d.Get("team").(string)

	log.Printf("[DEBUG] dataSourceOrganizationTeamMemberRead: %s", organization)

	cfg := &go_deploygate.ListOrganizationTeamMembersRequest{
		Organization: organization,
		Team:         team,
	}

	resp, err := client.ListOrganizationTeamMembers(cfg)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", organization, team))
	d.Set("users", resp.Users)

	return nil
}
