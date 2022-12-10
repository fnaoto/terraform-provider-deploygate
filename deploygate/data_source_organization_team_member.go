package deploygate

import (
	"context"
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganizationTeamMember() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves informantion about a existing organization member.",
		ReadContext: dataSourceOrganizationTeamMemberRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "Name of the organization. [Check your organizations](https://deploygate.com/organizations)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"team": {
				Description: "Name of the team in organization.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"users": {
				Description: "Data of the organization users.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of the user that is user or tester.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"name": {
							Description: "Name of the user",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"url": {
							Description: "Icon URL for user profile.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"icon_url": {
							Description: "URL of the user account.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationTeamMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", organization, team))
	d.Set("users", resp.Users)

	return nil
}
