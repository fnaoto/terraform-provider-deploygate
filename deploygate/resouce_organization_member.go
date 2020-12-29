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
	organization string
	members      []OrganizationMemberConfigMembers
}

// OrganizationMemberConfigMembers is config for OrganizationMemberConfig
type OrganizationMemberConfigMembers struct {
	Type     string
	name     string
	url      string
	iconURL  string
	inviting bool
}

func resourceOrganizationMemberRead(d *schema.ResourceData, meta interface{}) error {

	cfg := setOrganizationMemberConfig(d)

	log.Printf("[DEBUG] resourceOrganizationMemberRead")

	rs, _ := meta.(*Client).getOrganizationMember(cfg)

	d.SetId(cfg.organization)
	d.Set("members", rs.Members)

	return nil
}

func resourceOrganizationMemberCreate(d *schema.ResourceData, meta interface{}) error {
	cfg := setOrganizationMemberConfig(d)

	log.Printf("[DEBUG] resourceOrganizationMemberCreate")

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

	log.Printf("[DEBUG] resourceOrganizationMemberUpdate")

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

	log.Printf("[DEBUG] resourceOrganizationMemberDelete")

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
	for _, member := range cfg.members {
		g := &go_deploygate.AddOrganizationMemberInput{
			OrganizationName: cfg.organization,
			UserName:         member.name,
		}

		_, err := clt.client.AddOrganizationMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func (clt *Client) deleteOrganizationMember(cfg *OrganizationMemberConfig) error {
	for _, member := range cfg.members {
		g := &go_deploygate.DeleteOrganizationMemberInput{
			OrganizationName: cfg.organization,
			UserName:         member.name,
		}

		_, err := clt.client.DeleteOrganizationMember(g)

		if err != nil {
			return err
		}
	}
	return nil
}

func setOrganizationMemberConfig(d *schema.ResourceData) *OrganizationMemberConfig {
	var members []OrganizationMemberConfigMembers

	if v, ok := d.GetOk("members"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			members = append(members, OrganizationMemberConfigMembers{
				Type:     elem["type"].(string),
				name:     elem["name"].(string),
				url:      elem["url"].(string),
				iconURL:  elem["icon_url"].(string),
				inviting: elem["inviting"].(bool),
			})
		}
	}

	return &OrganizationMemberConfig{
		organization: d.Get("organization").(string),
		members:      members,
	}
}
