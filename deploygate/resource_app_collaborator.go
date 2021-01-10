package deploygate

import (
	"errors"
	"fmt"
	"log"
	"strings"

	go_deploygate "github.com/fnaoto/go-deploygate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppCollaborator() *schema.Resource {
	return &schema.Resource{
		Read:   resourceAppCollaboratorRead,
		Create: resourceAppCollaboratorCreate,
		Update: resourceAppCollaboratorUpdate,
		Delete: resourceAppCollaboratorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
	Users    []*go_deploygate.Collaborator
}

func resourceAppCollaboratorRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorRead")

	cfg, err := setAppCollaboratorConfig(d)

	if err != nil {
		return err
	}

	rs, gerr := meta.(*Client).getAppCollaborator(cfg)

	if gerr != nil {
		return gerr
	}

	var users []*go_deploygate.Collaborator

	for _, cus := range cfg.Users {
		for _, rus := range rs.Users {
			if cus.Name == rus.Name {
				users = append(users, rus)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", users)

	return nil
}

func resourceAppCollaboratorCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorCreate")

	cfg, err := setAppCollaboratorConfig(d)

	if err != nil {
		return err
	}

	aerr := meta.(*Client).addAppCollaborator(cfg)

	if aerr != nil {
		return aerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", cfg.Users)

	return nil
}

func resourceAppCollaboratorUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] resourceAppCollaboratorUpdate")

	cfg, err := setAppCollaboratorConfig(d)

	if err != nil {
		return err
	}

	derr := meta.(*Client).deleteAppCollaborator(cfg)

	if derr != nil {
		return derr
	}

	aerr := meta.(*Client).addAppCollaborator(cfg)

	if aerr != nil {
		return aerr
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cfg.Owner, cfg.Platform, cfg.AppID))
	d.Set("users", cfg.Users)

	return nil
}

func resourceAppCollaboratorDelete(d *schema.ResourceData, meta interface{}) error {
	cfg, err := setAppCollaboratorConfig(d)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] resourceAppCollaboratorDelete")

	derr := meta.(*Client).deleteAppCollaborator(cfg)

	if derr != nil {
		return derr
	}

	d.SetId("")

	return nil
}

func (clt *Client) getAppCollaborator(cfg *AppCollaboratorConfig) (*go_deploygate.GetAppCollaboratorResponseResult, error) {
	g := &go_deploygate.GetAppCollaboratorInput{
		Owner:    cfg.Owner,
		Platform: cfg.Platform,
		AppId:    cfg.AppID,
	}

	log.Printf("[DEBUG] getAppCollaborator: %s", g)

	collaborator, err := clt.client.GetAppCollaborator(g)

	if err != nil {
		return nil, err
	}

	return collaborator.Results, nil
}

func (clt *Client) addAppCollaborator(cfg *AppCollaboratorConfig) error {
	for _, user := range cfg.Users {
		g := &go_deploygate.AddAppCollaboratorInput{
			Owner:    cfg.Owner,
			Platform: cfg.Platform,
			AppId:    cfg.AppID,
			Users:    user.Name,
			Role:     int(user.Role),
		}

		_, err := clt.client.AddAppCollaborator(g)

		if err != nil {
			return err
		}

	}

	return nil
}

func (clt *Client) deleteAppCollaborator(cfg *AppCollaboratorConfig) error {
	for _, user := range cfg.Users {
		g := &go_deploygate.DeleteAppCollaboratorInput{
			Owner:    cfg.Owner,
			Platform: cfg.Platform,
			AppId:    cfg.AppID,
			Users:    user.Name,
		}

		_, err := clt.client.DeleteAppCollaborator(g)

		if err != nil {
			return err
		}
	}

	return nil
}

func setAppCollaboratorConfig(d *schema.ResourceData) (*AppCollaboratorConfig, error) {
	var users []*go_deploygate.Collaborator

	if v, ok := d.GetOk("users"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			users = append(users, &go_deploygate.Collaborator{
				Name: elem["name"].(string),
				Role: uint(elem["role"].(int)),
			})
		}
	}

	if !strings.Contains(d.Id(), "/") {
		return nil, errors.New("Not contains `/` in Id")
	}

	id := strings.Split(d.Id(), "/")

	if len(id) != 3 {
		return nil, errors.New("Not contains Owner, Platform or AppID in Id")
	}

	return &AppCollaboratorConfig{
		Owner:    id[0],
		Platform: id[1],
		AppID:    id[2],
		Users:    users,
	}, nil
}
