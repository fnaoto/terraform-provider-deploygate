package deploygate

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	go_deploygate "github.com/recruit-mp/go-deploygate"
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
						},
					},
				},
			},
			"teams": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"usage": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeInt,
			},
		},
	}
}

// AppCollaboratorConfig : is config for go-deploygate
type AppCollaboratorConfig struct {
	owner    string
	platform string
	appID    string
	users    string
}

func resourceAppCollaboratorRead(d *schema.ResourceData, meta interface{}) error {

	acc := setAppCollaboratorConfig(d)

	log.Printf("[DEBUG] resourceAppCollaboratorRead: %s", acc)

	rs, _ := meta.(*Client).getAppCollaborator(acc)

	d.SetId(fmt.Sprintf("%s/%s/%s", acc.owner, acc.platform, acc.appID))
	d.Set("users", rs.Users)
	d.Set("teams", rs.Teams)
	d.Set("usage", map[string]interface{}{
		"max":  rs.Usage.Max,
		"used": rs.Usage.Used,
	})

	return nil
}

func resourceAppCollaboratorCreate(d *schema.ResourceData, meta interface{}) error {
	acc := setAppCollaboratorConfig(d)

	log.Printf("[DEBUG] resourceAppCollaboratorCreate: %s", acc)

	rs, err := meta.(*Client).addAppCollaborator(acc)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", acc.owner, acc.platform, acc.appID))
	d.Set("users", rs.Added)

	return nil
}

func resourceAppCollaboratorUpdate(d *schema.ResourceData, meta interface{}) error {
	acc := setAppCollaboratorConfig(d)

	log.Printf("[DEBUG] resourceAppCollaboratorUpdate: %s", acc)

	_, err := meta.(*Client).deleteAppCollaborator(acc)

	if err != nil {
		return err
	}

	rs, _ := meta.(*Client).addAppCollaborator(acc)

	d.SetId(fmt.Sprintf("%s/%s/%s", acc.owner, acc.platform, acc.appID))
	d.Set("users", rs.Added)

	return nil
}

func resourceAppCollaboratorDelete(d *schema.ResourceData, meta interface{}) error {
	acc := setAppCollaboratorConfig(d)

	log.Printf("[DEBUG] resourceAppCollaboratorDelete: %s", acc)

	_, err := meta.(*Client).deleteAppCollaborator(acc)

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func (clt *Client) getAppCollaborator(cfg *AppCollaboratorConfig) (*go_deploygate.GetAppCollaboratorResponseResult, error) {

	g := &go_deploygate.GetAppCollaboratorInput{
		Owner:    cfg.owner,
		Platform: cfg.platform,
		AppId:    cfg.appID,
	}

	collaborator, err := clt.client.GetAppCollaborator(g)

	if err != nil {
		return nil, err
	}

	return collaborator.Results, nil
}

func (clt *Client) addAppCollaborator(cfg *AppCollaboratorConfig) (*go_deploygate.AddAppCollaboratorResponseResult, error) {
	g := &go_deploygate.AddAppCollaboratorInput{
		Owner:    cfg.owner,
		Platform: cfg.platform,
		AppId:    cfg.appID,
		Users:    cfg.users,
		Role:     1,
	}

	collaborator, err := clt.client.AddAppCollaborator(g)

	if err != nil {
		return nil, err
	}

	return collaborator.Results, nil
}

func (clt *Client) deleteAppCollaborator(cfg *AppCollaboratorConfig) (*go_deploygate.DeleteAppCollaboratorResponse, error) {
	g := &go_deploygate.DeleteAppCollaboratorInput{
		Owner:    cfg.owner,
		Platform: cfg.platform,
		AppId:    cfg.appID,
		Users:    cfg.users,
	}

	collaborator, err := clt.client.DeleteAppCollaborator(g)

	if err != nil {
		return nil, err
	}

	return collaborator, nil
}

func setAppCollaboratorConfig(d *schema.ResourceData) *AppCollaboratorConfig {
	var usersList string

	if v, ok := d.GetOk("users"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})
			usersList += elem["name"].(string) + ","
		}
	}

	acc := &AppCollaboratorConfig{
		owner:    d.Get("owner").(string),
		platform: d.Get("platform").(string),
		appID:    d.Get("app_id").(string),
		users:    usersList,
	}

	return acc
}
