package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnterpriseMember() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEnterpriseMemberRead,
		Schema: map[string]*schema.Schema{
			"enterprise": {
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
						"icon_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"last_access_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEnterpriseMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).client

	enterprise := d.Get("enterprise").(string)

	log.Printf("[DEBUG] dataSourceEnterpriseMemberRead: %s", enterprise)

	e := &go_deploygate.ListEnterpriseMembersRequest{
		Enterprise: enterprise,
	}

	rs, err := client.ListEnterpriseMembers(e)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s", enterprise))
	d.Set("users", rs.Users)

	return nil
}
