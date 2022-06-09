package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationMemberRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
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

func dataSourceOrganizationMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).client

	organization := d.Get("organization").(string)

	log.Printf("[DEBUG] dataSourceOrganizationMemberRead: %s", organization)

	g := &go_deploygate.ListOrganizationMembersRequest{
		Organization: organization,
	}

	om, err := client.ListOrganizationMembers(g)

	if err != nil {
		return err
	}

	rs := om

	d.SetId(fmt.Sprintf("%s", organization))
	d.Set("members", rs.Members)

	return nil
}
