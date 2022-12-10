package deploygate

import (
	"context"
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnterpriseOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves informantion about a existing enterprise organization member.",
		ReadContext: dataSourceEnterpriseOrganizationMemberRead,

		Schema: map[string]*schema.Schema{
			"enterprise": {
				Description: "Name of the enterprise. [Check your enterprises](https://deploygate.com/enterprises)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"organization": {
				Description: "Name of the organization in enterprise. [Check your enterprises](https://deploygate.com/enterprises)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"users": {
				Description: "Data of the enterprise users.",
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
						"icon_url": {
							Description: "Icon URL for user profile.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"url": {
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

func dataSourceEnterpriseOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Config).client

	enterprise := d.Get("enterprise").(string)
	organization := d.Get("organization").(string)

	log.Printf("[DEBUG] dataSourceEnterpriseOrganizationMemberRead: %s,%s", enterprise, organization)

	req := &go_deploygate.ListEnterpriseOrganizationMembersRequest{
		Enterprise:   enterprise,
		Organization: organization,
	}

	resp, err := client.ListEnterpriseOrganizationMembers(req)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", enterprise, organization))
	d.Set("users", resp.Users)

	return nil
}
