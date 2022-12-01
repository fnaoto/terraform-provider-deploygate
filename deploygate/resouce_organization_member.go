package deploygate

import (
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a organization member resource.",
		Read:        resourceOrganizationMemberRead,
		Create:      resourceOrganizationMemberCreate,
		Update:      resourceOrganizationMemberUpdate,
		Delete:      resourceOrganizationMemberDelete,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "Name of the organization. [Check your organizations](https://deploygate.com/organizations)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"members": {
				Description: "Data of the organization members.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of the user that is user or tester.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Name of the user",
							Type:        schema.TypeString,
							Required:    true,
						},
						"icon_url": {
							Description: "Icon URL for user profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"url": {
							Description: "URL of the user account.",
							Type:        schema.TypeString,
							Computed:    true,
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
	log.Printf("[DEBUG] resourceOrganizationMemberRead %#v", d)

	cfg := setOrganizationMemberConfig(d)
	resp, err := meta.(*Config).getOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Organization)
	d.Set("members", resp.Members)

	return nil
}

func resourceOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberCreate %#v", d)

	cfg := setOrganizationMemberConfig(d)
	err := meta.(*Config).addOrganizationMember(cfg)

	if err != nil {
		return err
	}

	resp, err := meta.(*Config).getOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Organization)
	d.Set("members", resp.Members)

	return nil
}

func resourceOrganizationMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberUpdate %#v", d)

	var err error

	cfg := setOrganizationMemberConfig(d)
	err = meta.(*Config).deleteOrganizationMember(cfg)

	if err != nil {
		return err
	}

	err = meta.(*Config).addOrganizationMember(cfg)

	if err != nil {
		return err
	}

	resp, err := meta.(*Config).getOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Organization)
	d.Set("members", resp.Members)

	return nil
}

func resourceOrganizationMemberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberDelete %#v", d)

	cfg := setOrganizationMemberConfig(d)
	err := meta.(*Config).deleteOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func (c *Config) getOrganizationMember(cfg *OrganizationMemberConfig) (*go_deploygate.ListOrganizationMembersResponse, error) {
	log.Printf("[DEBUG] getOrganizationMember: %#v", cfg)

	req := &go_deploygate.ListOrganizationMembersRequest{
		Organization: cfg.Organization,
	}

	resp, err := c.client.ListOrganizationMembers(req)

	if err != nil {
		return nil, err
	}

	var members []go_deploygate.Member

	for _, csm := range cfg.Members {
		for _, rsm := range resp.Members {
			if csm.Name == rsm.Name {
				members = append(members, rsm)
			}
		}
	}

	resp.Members = members

	return resp, nil
}

func (c *Config) addOrganizationMember(cfg *OrganizationMemberConfig) error {
	log.Printf("[DEBUG] addOrganizationMember: %#v", cfg)

	for _, member := range cfg.Members {
		req := &go_deploygate.AddOrganizationMemberByUserNameRequest{
			Organization: cfg.Organization,
			UserName:     member.Name,
		}

		_, err := c.client.AddOrganizationMemberByUserName(req)

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) deleteOrganizationMember(cfg *OrganizationMemberConfig) error {
	log.Printf("[DEBUG] deleteOrganizationMember %#v", cfg)

	for _, member := range cfg.Members {
		g := &go_deploygate.RemoveOrganizationMemberByUserNameRequest{
			Organization: cfg.Organization,
			UserName:     member.Name,
		}

		_, err := c.client.RemoveOrganizationMemberByUserName(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func setOrganizationMemberConfig(d *schema.ResourceData) *OrganizationMemberConfig {
	log.Printf("[DEBUG] setOrganizationMemberConfig %#v", d)

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
