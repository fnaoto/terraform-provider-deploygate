package deploygate

import (
	"context"
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationTeamMember() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a organization member resource.",
		ReadContext:   resourceOrganizationTeamMemberRead,
		CreateContext: resourceOrganizationTeamMemberCreate,
		UpdateContext: resourceOrganizationTeamMemberUpdate,
		DeleteContext: resourceOrganizationTeamMemberDelete,

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
				Description: "Data of the organization team users.",
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

// OrganizationTeamMemberConfig is config for go-deploygate
type OrganizationTeamMemberConfig struct {
	Organization string
	Team         string
	Users        []*go_deploygate.Member
}

func resourceOrganizationTeamMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceOrganizationTeamMemberRead %#v", d)

	cfg := setOrganizationTeamMemberConfig(d)
	resp, err := meta.(*Config).getOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cfg.Organization, cfg.Team))
	d.Set("users", resp.Users)

	return nil
}

func resourceOrganizationTeamMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceOrganizationTeamMemberCreate %#v", d)

	cfg := setOrganizationTeamMemberConfig(d)
	err := meta.(*Config).addOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := meta.(*Config).getOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cfg.Organization, cfg.Team))
	d.Set("users", resp.Users)

	return nil
}

func resourceOrganizationTeamMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceOrganizationTeamMemberUpdate %#v", d)

	var err error

	cfg := setOrganizationTeamMemberConfig(d)
	err = meta.(*Config).deleteOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	err = meta.(*Config).addOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := meta.(*Config).getOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cfg.Organization, cfg.Team))
	d.Set("users", resp.Users)

	return nil
}

func resourceOrganizationTeamMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] resourceOrganizationTeamMemberDelete %#v", d)

	cfg := setOrganizationTeamMemberConfig(d)
	err := meta.(*Config).deleteOrganizationTeamMember(cfg)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func (c *Config) getOrganizationTeamMember(cfg *OrganizationTeamMemberConfig) (*go_deploygate.ListOrganizationTeamMembersResponse, error) {
	log.Printf("[DEBUG] getOrganizationTeamMember: %#v", cfg)

	req := &go_deploygate.ListOrganizationTeamMembersRequest{
		Organization: cfg.Organization,
		Team:         cfg.Team,
	}

	resp, err := c.client.ListOrganizationTeamMembers(req)

	if err != nil {
		return nil, err
	}

	var users []go_deploygate.Member

	for _, csm := range cfg.Users {
		for _, rsm := range resp.Users {
			if csm.Name == rsm.Name {
				users = append(users, rsm)
			}
		}
	}

	resp.Users = users

	return resp, nil
}

func (c *Config) addOrganizationTeamMember(cfg *OrganizationTeamMemberConfig) error {
	log.Printf("[DEBUG] addOrganizationTeamMember %#v", cfg)

	for _, user := range cfg.Users {
		req := &go_deploygate.AddOrganizationTeamMemberRequest{
			Organization: cfg.Organization,
			Team:         cfg.Team,
			User:         user.Name,
		}

		_, err := c.client.AddOrganizationTeamMember(req)

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) deleteOrganizationTeamMember(cfg *OrganizationTeamMemberConfig) error {
	log.Printf("[DEBUG] deleteOrganizationTeamMember %#v", cfg)

	for _, user := range cfg.Users {
		req := &go_deploygate.RemoveOrganizationTeamMemberRequest{
			Organization: cfg.Organization,
			Team:         cfg.Team,
			User:         user.Name,
		}

		_, err := c.client.RemoveOrganizationTeamMember(req)

		if err != nil {
			return err
		}
	}
	return nil
}

func setOrganizationTeamMemberConfig(d *schema.ResourceData) *OrganizationTeamMemberConfig {
	log.Printf("[DEBUG] setOrganizationTeamMemberConfig %#v", d)

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

	return &OrganizationTeamMemberConfig{
		Organization: d.Get("organization").(string),
		Team:         d.Get("team").(string),
		Users:        users,
	}
}
