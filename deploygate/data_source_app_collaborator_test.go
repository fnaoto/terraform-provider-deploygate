package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_DataSourceAppCollaborator_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { Test_DGPreCheck(t) },
		Providers: testDGProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAppCollaboratorConfig,
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAppCollaborator("data.deploygate_app_collaborator.current"),
				),
			},
		},
	})
}

func testDataSourceAppCollaborator(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}

		if rs.Primary.Attributes["owner"] == "" {
			return fmt.Errorf("owner expected to not be nil")
		}

		if rs.Primary.Attributes["platform"] == "" {
			return fmt.Errorf("platform expected to not be nil")
		}

		if rs.Primary.Attributes["app_id"] == "" {
			return fmt.Errorf("app_id expected to not be nil")
		}

		return nil
	}
}

const testDataSourceAppCollaboratorConfig = `
provider "deploygate" {
	api_key = var.user_api_key
}

variable "user_api_key" {
  type = string
}

data "deploygate_app_collaborator" "current" {
	owner    = var.owner
	platform = var.platform
	app_id   = var.app_id
}

variable "platform" {
  type = string
}

variable "app_id" {
  type = string
}

variable "owner" {
  type = string
}
`
