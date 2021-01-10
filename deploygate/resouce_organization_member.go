package deploygate

import (
	"log"

	go_deploygate "github.com/fnaoto/go-deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Read:   resourceOrganizationMemberRead,
		Create: resourceOrganizationMemberCreate,
		Update: resourceOrganizationMemberUpdate,
		Delete: resourceOrganizationMemberDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
						"inviting": {
							Type:     schema.TypeBool,
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

	d.SetId(cfg.Organization)
	d.Set("members", cfg.Members)

	return nil
}

func resourceOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceOrganizationMemberCreate")

	cfg := setOrganizationMemberConfig(d)
	err := meta.(*Client).addOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId(cfg.Organization)
	d.Set("members", cfg.Members)

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

	d.SetId(cfg.Organization)
	d.Set("members", cfg.Members)

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

func (clt *Client) getOrganizationMember(cfg *OrganizationMemberConfig) (*go_deploygate.GetOrganizationMemberResponse, error) {
	g := &go_deploygate.GetOrganizationMemberInput{
		OrganizationName: cfg.Organization,
	}

	rs, err := clt.client.GetOrganizationMember(g)

	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (clt *Client) addOrganizationMember(cfg *OrganizationMemberConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.AddOrganizationMemberInput{
			OrganizationName: cfg.Organization,
			UserName:         member.Name,
		}

		_, err := clt.client.AddOrganizationMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func (clt *Client) deleteOrganizationMember(cfg *OrganizationMemberConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.DeleteOrganizationMemberInput{
			OrganizationName: cfg.Organization,
			UserName:         member.Name,
		}

		_, err := clt.client.DeleteOrganizationMember(g)

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
				Type:     elem["type"].(string),
				Name:     elem["name"].(string),
				URL:      elem["url"].(string),
				IconURL:  elem["icon_url"].(string),
				Inviting: elem["inviting"].(bool),
			})
		}
	}

	return &OrganizationMemberConfig{
		Organization: d.Get("organization").(string),
		Members:      members,
	}
}
