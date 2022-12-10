package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnterpriseOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a enterprise organization member resource.",

		Read:   resourceEnterpriseOrganizationMemberRead,
		Create: resourceEnterpriseOrganizationMemberCreate,
		Update: resourceEnterpriseOrganizationMemberUpdate,
		Delete: resourceEnterpriseOrganizationMemberDelete,

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

// EnterpriseOrganizationMemberConfig is config for go-deploygate
type EnterpriseOrganizationMemberConfig struct {
	Enterprise   string
	Organization string
	Users        []*go_deploygate.Member
}

func resourceEnterpriseOrganizationMemberRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseOrganizationMemberRead")

	cfg := setEnterpriseOrganizationMemberConfig(d)
	resp, err := meta.(*Config).listEnterpriseOrganizationMembers(cfg)

	if err != nil {
		return err
	}

	users := resp.Users

	d.SetId(fmt.Sprintf("%s/%s", cfg.Enterprise, cfg.Organization))
	d.Set("users", users)

	return nil
}

func resourceEnterpriseOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseOrganizationMemberCreate")

	cfg := setEnterpriseOrganizationMemberConfig(d)
	err := meta.(*Config).addEnterpriseOrganizationMember(cfg)

	if err != nil {
		return err
	}

	resp, err := meta.(*Config).listEnterpriseOrganizationMembers(cfg)

	if err != nil {
		return err
	}

	users := resp.Users

	d.SetId(fmt.Sprintf("%s/%s", cfg.Enterprise, cfg.Organization))
	d.Set("users", users)

	return nil
}

func resourceEnterpriseOrganizationMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseOrganizationMemberUpdate")

	cfg := setEnterpriseOrganizationMemberConfig(d)
	err := meta.(*Config).deleteEnterpriseOrganizationMember(cfg)

	if err != nil {
		return err
	}

	err = meta.(*Config).addEnterpriseOrganizationMember(cfg)

	if err != nil {
		return err
	}

	resp, err := meta.(*Config).listEnterpriseOrganizationMembers(cfg)

	if err != nil {
		return err
	}

	users := resp.Users

	d.SetId(fmt.Sprintf("%s/%s", cfg.Enterprise, cfg.Organization))
	d.Set("users", users)

	return nil
}

func resourceEnterpriseOrganizationMemberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseOrganizationMemberDelete")

	cfg := setEnterpriseOrganizationMemberConfig(d)
	err := meta.(*Config).deleteEnterpriseOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func (c *Config) listEnterpriseOrganizationMembers(cfg *EnterpriseOrganizationMemberConfig) (*go_deploygate.ListEnterpriseOrganizationMembersResponse, error) {
	req := &go_deploygate.ListEnterpriseOrganizationMembersRequest{
		Enterprise:   cfg.Enterprise,
		Organization: cfg.Organization,
	}

	log.Printf("[DEBUG] listEnterpriseOrganizationMembers: %s", cfg.Enterprise)

	resp, err := c.client.ListEnterpriseOrganizationMembers(req)

	if err != nil {
		return nil, err
	}

	var users []go_deploygate.Member

	for _, c := range cfg.Users {
		for _, user := range resp.Users {
			if c.Name == user.Name {
				users = append(users, user)
			}
		}
	}

	resp.Users = users

	return resp, nil
}

func (c *Config) addEnterpriseOrganizationMember(cfg *EnterpriseOrganizationMemberConfig) error {
	for _, user := range cfg.Users {
		g := &go_deploygate.AddEnterpriseOrganizationMemberRequest{
			Enterprise:   cfg.Enterprise,
			Organization: cfg.Organization,
			User:         user.Name,
		}

		_, err := c.client.AddEnterpriseOrganizationMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) deleteEnterpriseOrganizationMember(cfg *EnterpriseOrganizationMemberConfig) error {
	for _, user := range cfg.Users {
		req := &go_deploygate.RemoveEnterpriseOrganizationMemberRequest{
			Enterprise:   cfg.Enterprise,
			Organization: cfg.Organization,
			User:         user.Name,
		}

		_, err := c.client.RemoveEnterpriseOrganizationMember(req)

		if err != nil {
			return err
		}
	}
	return nil
}

func setEnterpriseOrganizationMemberConfig(d *schema.ResourceData) *EnterpriseOrganizationMemberConfig {
	var users []*go_deploygate.Member

	if v, ok := d.GetOk("users"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			users = append(users, &go_deploygate.Member{
				Type:    elem["type"].(string),
				Name:    elem["name"].(string),
				IconUrl: elem["icon_url"].(string),
				Url:     elem["url"].(string),
			})
		}
	}

	return &EnterpriseOrganizationMemberConfig{
		Enterprise:   d.Get("enterprise").(string),
		Organization: d.Get("organization").(string),
		Users:        users,
	}
}
