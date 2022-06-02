package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppCollaborator() *schema.Resource {
	return &schema.Resource{
		Read:   resourceAppCollaboratorRead,
		Create: resourceAppCollaboratorCreate,
		Update: resourceAppCollaboratorUpdate,
		Delete: resourceAppCollaboratorDelete,

		Schema: map[string]*schema.Schema{
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
		},
	}
}

// AppCollaboratorConfig is config for go-deploygate
type AppCollaboratorConfig struct {
	Owner    string
	Platform string
	AppID    string
	Members  []*go_deploygate.Member
}

func resourceAppCollaboratorRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorRead")

	cfg := setAppCollaboratorConfig(d)
	rs, gerr := meta.(*Client).getAppCollaborator(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", rs.Users)

	return nil
}

func resourceAppCollaboratorCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorCreate")

	cfg := setAppCollaboratorConfig(d)
	aerr := meta.(*Client).addAppCollaborator(cfg)

	if aerr != nil {
		return aerr
	}

	rs, gerr := meta.(*Client).getAppCollaborator(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", rs.Users)

	return nil
}

func resourceAppCollaboratorUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorUpdate")

	cfg := setAppCollaboratorConfig(d)
	derr := meta.(*Client).deleteAppCollaborator(cfg)

	if derr != nil {
		return derr
	}

	aerr := meta.(*Client).addAppCollaborator(cfg)

	if aerr != nil {
		return aerr
	}

	rs, gerr := meta.(*Client).getAppCollaborator(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", rs.Users)

	return nil
}

func resourceAppCollaboratorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorDelete")

	cfg := setAppCollaboratorConfig(d)
	derr := meta.(*Client).deleteAppCollaborator(cfg)

	if derr != nil {
		return derr
	}

	d.SetId("")

	return nil
}

func (clt *Client) getAppCollaborator(cfg *AppCollaboratorConfig) (*go_deploygate.GetAppMembersResponseResult, error) {
	g := &go_deploygate.GetAppMembersRequest{
		Owner:    cfg.Owner,
		Platform: cfg.Platform,
		AppId:    cfg.AppID,
	}

	log.Printf("[DEBUG] getAppCollaborator: %s", g)

	rs, err := clt.client.GetAppMembers(g)

	if err != nil {
		return nil, err
	}

	var members []*go_deploygate.Member

	for _, cus := range cfg.Members {
		for _, rus := range rs.Results.Users {
			if cus.Name == rus.Name {
				members = append(members, rus)
			}
		}
	}

	rs.Results.Users = members

	return rs.Results, nil
}

func (clt *Client) addAppCollaborator(cfg *AppCollaboratorConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.AddAppMembersRequest{
			Owner:    cfg.Owner,
			Platform: cfg.Platform,
			AppId:    cfg.AppID,
			Users:    member.Name,
			Role:     fmt.Sprint(member.Role),
		}

		_, err := clt.client.AddAppMembers(g)

		if err != nil {
			return err
		}

	}

	return nil
}

func (clt *Client) deleteAppCollaborator(cfg *AppCollaboratorConfig) error {
	for _, member := range cfg.Members {
		g := &go_deploygate.RemoveAppMembersRequest{
			Owner:    cfg.Owner,
			Platform: cfg.Platform,
			AppId:    cfg.AppID,
			Users:    member.Name,
		}

		_, err := clt.client.RemoveAppMembers(g)

		if err != nil {
			return err
		}
	}

	return nil
}

func setAppCollaboratorConfig(d *schema.ResourceData) *AppCollaboratorConfig {
	var members []*go_deploygate.Member

	if v, ok := d.GetOk("users"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			members = append(members, &go_deploygate.Member{
				Name: elem["name"].(string),
				Role: uint(elem["role"].(int)),
			})
		}
	}

	return &AppCollaboratorConfig{
		Owner:    d.Get("owner").(string),
		Platform: d.Get("platform").(string),
		AppID:    d.Get("app_id").(string),
		Members:  members,
	}
}
