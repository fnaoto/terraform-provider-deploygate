package deploygate

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	go_deploygate "github.com/recruit-mp/go-deploygate"
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
						"inviting": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// OrganizationMemberConfig is config for go-deploygate
type OrganizationMemberConfig struct {
	organization string
}

func resourceOrganizationMemberRead(d *schema.ResourceData, meta interface{}) error {

	cfg := setOrganizationMemberConfig(d)

	log.Printf("[DEBUG] resourceOrganizationMemberRead: %s", cfg)

	rs, _ := meta.(*Client).getOrganizationMember(cfg)

	d.SetId(cfg.organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	cfg := setOrganizationMemberConfig(d)

	log.Printf("[DEBUG] resourceOrganizationMemberCreate %s", cfg)

	err := meta.(*Client).addOrganizationMember(cfg)

	if err != nil {
		return err
	}

	rs, _ := meta.(*Client).getOrganizationMember(cfg)

	d.SetId(cfg.organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	cfg := setOrganizationMemberConfig(d)

	log.Printf("[DEBUG] resourceOrganizationMemberUpdate %s", cfg)

	derr := meta.(*Client).deleteOrganizationMember(cfg)

	if derr != nil {
		return derr
	}

	aerr := meta.(*Client).addOrganizationMember(cfg)

	if aerr != nil {
		return aerr
	}

	rs, _ := meta.(*Client).getOrganizationMember(cfg)

	d.SetId(cfg.organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberDelete(d *schema.ResourceData, meta interface{}) error {
	cfg := setOrganizationMemberConfig(d)

	log.Printf("[DEBUG] resourceOrganizationMemberDelete %s", cfg)

	err := meta.(*Client).deleteOrganizationMember(cfg)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func (clt *Client) getOrganizationMember(cfg *OrganizationMemberConfig) (*go_deploygate.GetOrganizationMemberResponse, error) {

	g := &go_deploygate.GetOrganizationMemberInput{
		OrganizationName: cfg.organization,
	}

	rs, err := clt.client.GetOrganizationMember(g)

	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (clt *Client) addOrganizationMember(cfg *OrganizationMemberConfig) error {
	g := &go_deploygate.AddOrganizationMemberInput{
		OrganizationName: cfg.organization,
	}

	_, err := clt.client.AddOrganizationMember(g)

	if err != nil {
		return err
	}

	return nil
}

func (clt *Client) deleteOrganizationMember(cfg *OrganizationMemberConfig) error {
	g := &go_deploygate.DeleteOrganizationMemberInput{
		OrganizationName: cfg.organization,
	}

	_, err := clt.client.DeleteOrganizationMember(g)

	if err != nil {
		return err
	}

	return nil
}

func setOrganizationMemberConfig(d *schema.ResourceData) *OrganizationMemberConfig {
	return &OrganizationMemberConfig{
		organization: d.Get("organization").(string),
	}
}
