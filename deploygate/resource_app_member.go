package deploygate

import (
	"fmt"
	"log"

	go_deploygate "github.com/fnaoto/go_deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppMember() *schema.Resource {
	return &schema.Resource{
		Read:   resourceAppMemberRead,
		Create: resourceAppMemberCreate,
		Update: resourceAppMemberUpdate,
		Delete: resourceAppMemberDelete,

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
							Default:  2,
						},
					},
				},
			},
		},
	}
}

// AppMemberConfig is config for go-deploygate
type AppMemberConfig struct {
	Owner    string
	Platform string
	AppID    string
	Users    []*go_deploygate.User
}

func resourceAppMemberRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppMemberRead")

	cfg := setAppMemberConfig(d)
	rs, gerr := meta.(*Config).getAppMember(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", rs.Users)

	return nil
}

func resourceAppMemberCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppMemberCreate")

	cfg := setAppMemberConfig(d)
	aerr := meta.(*Config).addAppMember(cfg)

	if aerr != nil {
		return aerr
	}

	rs, gerr := meta.(*Config).getAppMember(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", rs.Users)

	return nil
}

func resourceAppMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppMemberUpdate")

	cfg := setAppMemberConfig(d)
	derr := meta.(*Config).deleteAppMember(cfg)

	if derr != nil {
		return derr
	}

	aerr := meta.(*Config).addAppMember(cfg)

	if aerr != nil {
		return aerr
	}

	rs, gerr := meta.(*Config).getAppMember(cfg)

	if gerr != nil {
		return gerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", rs.Users)

	return nil
}

func resourceAppMemberDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppMemberDelete")

	cfg := setAppMemberConfig(d)
	derr := meta.(*Config).deleteAppMember(cfg)

	if derr != nil {
		return derr
	}

	d.SetId("")

	return nil
}

func (c *Config) getAppMember(cfg *AppMemberConfig) (*go_deploygate.GetAppMembersResponseResult, error) {
	g := &go_deploygate.GetAppMembersRequest{
		Owner:    cfg.Owner,
		Platform: cfg.Platform,
		AppId:    cfg.AppID,
	}

	log.Printf("[DEBUG] getAppMember: %s", g)

	rs, err := c.client.GetAppMembers(g)

	if err != nil {
		return nil, err
	}

	var users []go_deploygate.User

	for _, cus := range cfg.Users {
		for _, rus := range rs.Results.Users {
			if cus.Name == rus.Name {
				users = append(users, rus)
			}
		}
	}

	rs.Results.Users = users

	return &rs.Results, nil
}

func (c *Config) addAppMember(cfg *AppMemberConfig) error {
	for _, user := range cfg.Users {
		g := &go_deploygate.AddAppMembersRequest{
			Owner:    cfg.Owner,
			Platform: cfg.Platform,
			AppId:    cfg.AppID,
			Users:    user.Name,
			Role:     fmt.Sprint(user.Role),
		}

		_, err := c.client.AddAppMembers(g)

		if err != nil {
			return err
		}

	}

	return nil
}

func (c *Config) deleteAppMember(cfg *AppMemberConfig) error {
	for _, user := range cfg.Users {
		g := &go_deploygate.RemoveAppMembersRequest{
			Owner:    cfg.Owner,
			Platform: cfg.Platform,
			AppId:    cfg.AppID,
			Users:    user.Name,
		}

		_, err := c.client.RemoveAppMembers(g)

		if err != nil {
			return err
		}
	}

	return nil
}

func setAppMemberConfig(d *schema.ResourceData) *AppMemberConfig {
	var users []*go_deploygate.User

	if v, ok := d.GetOk("users"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			users = append(users, &go_deploygate.User{
				Name: elem["name"].(string),
				Role: uint(elem["role"].(int)),
			})
		}
	}

	return &AppMemberConfig{
		Owner:    d.Get("owner").(string),
		Platform: d.Get("platform").(string),
		AppID:    d.Get("app_id").(string),
		Users:    users,
	}
}
