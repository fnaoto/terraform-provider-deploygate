package deploygate

import (
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOrganizationMemberRead,
		Create: resourceOrganizationMemberCreate,
		Update: resourceOrganizationMemberUpdate,
		Delete: resourceOrganizationMemberDelete,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// OrganizationMemberConfig is config for go-deploygate
type OrganizationMemberConfig struct {
	Organization string
	Members      []*go_deploygate.Member
}

func resourceOrganizationMemberRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberRead")

	cfg := setOrganizationMemberConfig(d)
	rs, err := meta.(*Client).getOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberCreate")

	cfg := setOrganizationMemberConfig(d)
	aerr := meta.(*Client).addOrganizationMember(cfg)

	if aerr != nil {
		return aerr
	}

	rs, gerr := meta.(*Client).getOrganizationMember(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(cfg.Organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberUpdate")

	cfg := setOrganizationMemberConfig(d)
	derr := meta.(*Client).deleteOrganizationMember(cfg)

	if derr != nil {
		return derr
	}

	aerr := meta.(*Client).addOrganizationMember(cfg)

	if aerr != nil {
		return aerr
	}

	rs, gerr := meta.(*Client).getOrganizationMember(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(cfg.Organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberDelete")

	cfg := setOrganizationMemberConfig(d)
	err := meta.(*Client).deleteOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func (clt *Client) getOrganizationMember(cfg *OrganizationMemberConfig) (*go_deploygate.ListOrganizationMembersResponse, error) {
	g := &go_deploygate.ListOrganizationMembersRequest{
		Organization: cfg.Organization,
	}

	log.Printf("[DEBUG] getOrganizationMember: %s", cfg.Organization)

	rs, err := clt.client.ListOrganizationMembers(g)

	if err != nil {
		return nil, err
	}

	var members []go_deploygate.Member

	for _, csm := range cfg.Members {
		for _, rsm := range rs.Members {
			if csm.Name == rsm.Name {
				members = append(members, rsm)
			}
		}
	}

	rs.Members = members

	return rs, nil
}

func (clt *Client) addOrganizationMember(cfg *OrganizationMemberConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.AddOrganizationMemberByUserNameRequest{
			Organization: cfg.Organization,
			UserName:     member.Name,
		}

		_, err := clt.client.AddOrganizationMemberByUserName(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func (clt *Client) deleteOrganizationMember(cfg *OrganizationMemberConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.RemoveOrganizationMemberByUserNameRequest{
			Organization: cfg.Organization,
			UserName:     member.Name,
		}

		_, err := clt.client.RemoveOrganizationMemberByUserName(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func setOrganizationMemberConfig(d *schema.ResourceData) *OrganizationMemberConfig {
	var members []*go_deploygate.Member

	if v, ok := d.GetOk("members"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			members = append(members, &go_deploygate.Member{
				Type:    elem["type"].(string),
				Name:    elem["name"].(string),
				IconUrl: elem["icon_url"].(string),
				Url:     elem["url"].(string),
			})
		}
	}

	return &OrganizationMemberConfig{
		Organization: d.Get("organization").(string),
		Members:      members,
	}
}
