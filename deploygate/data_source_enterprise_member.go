package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnterpriseMember() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves informantion about a existing enterprise member.",
		Read:        dataSourceEnterpriseMemberRead,
		Schema: map[string]*schema.Schema{
			"enterprise": {
				Description: "Name of the enterprise. [Check your enterprises](https://deploygate.com/enterprises)",
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

func dataSourceEnterpriseMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).client

	enterprise := d.Get("enterprise").(string)

	log.Printf("[DEBUG] dataSourceEnterpriseMemberRead: %s", enterprise)

	req := &go_deploygate.ListEnterpriseMembersRequest{
		Enterprise: enterprise,
	}

	resp, err := client.ListEnterpriseMembers(req)

	if err != nil {
		return err
	}

	users := converEnterpriseMemberToMember(resp.Users)

	d.SetId(fmt.Sprintf("%s", enterprise))
	d.Set("users", users)

	return nil
}
