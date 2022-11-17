package deploygate

import (
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnterpriseMember() *schema.Resource {
	return &schema.Resource{
		Read:   resourceEnterpriseMemberRead,
		Create: resourceEnterpriseMemberCreate,
		Update: resourceEnterpriseMemberUpdate,
		Delete: resourceEnterpriseMemberDelete,

		Schema: map[string]*schema.Schema{
			"enterprise": {
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

// EnterpriseMemberConfig is config for go-deploygate
type EnterpriseMemberConfig struct {
	Enterprise string
	Members    []*go_deploygate.EnterpriseMember
}

func resourceEnterpriseMemberRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseMemberRead")

	cfg := setEnterpriseMemberConfig(d)
	rs, err := meta.(*Config).listEnterpriseMembers(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Enterprise)
	d.Set("members", rs.Users)

	return nil
}

func resourceEnterpriseMemberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseMemberCreate")

	cfg := setEnterpriseMemberConfig(d)
	err := meta.(*Config).addEnterpriseMember(cfg)

	if err != nil {
		return err
	}

	rs, err := meta.(*Config).listEnterpriseMembers(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Enterprise)
	d.Set("members", rs.Users)

	return nil
}

func resourceEnterpriseMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseMemberUpdate")

	cfg := setEnterpriseMemberConfig(d)
	err := meta.(*Config).deleteEnterpriseMember(cfg)

	if err != nil {
		return err
	}

	err = meta.(*Config).addEnterpriseMember(cfg)

	if err != nil {
		return err
	}

	rs, err := meta.(*Config).listEnterpriseMembers(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Enterprise)
	d.Set("members", rs.Users)

	return nil
}

func resourceEnterpriseMemberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceEnterpriseMemberDelete")

	cfg := setEnterpriseMemberConfig(d)
	err := meta.(*Config).deleteEnterpriseMember(cfg)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func (c *Config) listEnterpriseMembers(cfg *EnterpriseMemberConfig) (*go_deploygate.ListEnterpriseMembersResponse, error) {
	g := &go_deploygate.ListEnterpriseMembersRequest{
		Enterprise: cfg.Enterprise,
	}

	log.Printf("[DEBUG] listEnterpriseMembers: %s", cfg.Enterprise)

	rs, err := c.client.ListEnterpriseMembers(g)

	if err != nil {
		return nil, err
	}

	var members []go_deploygate.EnterpriseMember

	for _, c := range cfg.Members {
		for _, u := range rs.Users {
			if c.Name == u.Name {
				members = append(members, u)
			}
		}
	}

	rs.Users = members

	return rs, nil
}

func (c *Config) addEnterpriseMember(cfg *EnterpriseMemberConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.AddEnterpriseMemberRequest{
			Enterprise: cfg.Enterprise,
			User:       member.Name,
		}

		_, err := c.client.AddEnterpriseMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) deleteEnterpriseMember(cfg *EnterpriseMemberConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.RemoveEnterpriseMemberRequest{
			Enterprise: cfg.Enterprise,
			User:       member.Name,
		}

		_, err := c.client.RemoveEnterpriseMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func setEnterpriseMemberConfig(d *schema.ResourceData) *EnterpriseMemberConfig {
	var members []*go_deploygate.EnterpriseMember

	if v, ok := d.GetOk("members"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			members = append(members, &go_deploygate.EnterpriseMember{
				Type:    elem["type"].(string),
				Name:    elem["name"].(string),
				IconUrl: elem["icon_url"].(string),
				Url:     elem["url"].(string),
			})
		}
	}

	return &EnterpriseMemberConfig{
		Enterprise: d.Get("enterprise").(string),
		Members:    members,
	}
}
