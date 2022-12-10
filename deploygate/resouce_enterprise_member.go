package deploygate

import (
	"context"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnterpriseMember() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a enterprise member resource.",
		ReadContext:   resourceEnterpriseMemberRead,
		CreateContext: resourceEnterpriseMemberCreate,
		UpdateContext: resourceEnterpriseMemberUpdate,
		DeleteContext: resourceEnterpriseMemberDelete,

		Schema: map[string]*schema.Schema{
			"enterprise": {
				Description: "Name of the enterprise. [Check your enterprises](https://deploygate.com/enterprises)",
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

// EnterpriseMemberConfig is config for go-deploygate
type EnterpriseMemberConfig struct {
	Enterprise string
	Users      []*go_deploygate.Member
}

func resourceEnterpriseMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceEnterpriseMemberRead")

	cfg := setEnterpriseMemberConfig(d)
	resp, err := meta.(*Config).listEnterpriseMembers(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	users := converEnterpriseMemberToMember(resp.Users)

	d.SetId(cfg.Enterprise)
	d.Set("users", users)

	return nil
}

func resourceEnterpriseMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceEnterpriseMemberCreate")

	cfg := setEnterpriseMemberConfig(d)
	err := meta.(*Config).addEnterpriseMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := meta.(*Config).listEnterpriseMembers(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	users := converEnterpriseMemberToMember(resp.Users)

	d.SetId(cfg.Enterprise)
	d.Set("users", users)

	return nil
}

func resourceEnterpriseMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceEnterpriseMemberUpdate")

	cfg := setEnterpriseMemberConfig(d)
	err := meta.(*Config).deleteEnterpriseMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	err = meta.(*Config).addEnterpriseMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := meta.(*Config).listEnterpriseMembers(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	users := converEnterpriseMemberToMember(resp.Users)

	d.SetId(cfg.Enterprise)
	d.Set("users", users)

	return nil
}

func resourceEnterpriseMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceEnterpriseMemberDelete")

	cfg := setEnterpriseMemberConfig(d)
	err := meta.(*Config).deleteEnterpriseMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func (c *Config) listEnterpriseMembers(cfg *EnterpriseMemberConfig) (*go_deploygate.ListEnterpriseMembersResponse, error) {
	req := &go_deploygate.ListEnterpriseMembersRequest{
		Enterprise: cfg.Enterprise,
	}

	log.Printf("[DEBUG] listEnterpriseMembers: %s", cfg.Enterprise)

	resp, err := c.client.ListEnterpriseMembers(req)

	if err != nil {
		return nil, err
	}

	var users []go_deploygate.EnterpriseMember

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

func (c *Config) addEnterpriseMember(cfg *EnterpriseMemberConfig) error {
	for _, user := range cfg.Users {
		g := &go_deploygate.AddEnterpriseMemberRequest{
			Enterprise: cfg.Enterprise,
			User:       user.Name,
		}

		_, err := c.client.AddEnterpriseMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) deleteEnterpriseMember(cfg *EnterpriseMemberConfig) error {
	for _, user := range cfg.Users {
		req := &go_deploygate.RemoveEnterpriseMemberRequest{
			Enterprise: cfg.Enterprise,
			User:       user.Name,
		}

		_, err := c.client.RemoveEnterpriseMember(req)

		if err != nil {
			return err
		}
	}
	return nil
}

func setEnterpriseMemberConfig(d *schema.ResourceData) *EnterpriseMemberConfig {
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

	return &EnterpriseMemberConfig{
		Enterprise: d.Get("enterprise").(string),
		Users:      users,
	}
}
