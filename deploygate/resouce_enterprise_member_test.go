package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_ResourceEnterpriseMember_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { Test_DGPreCheck(t) },
		Providers: testDGProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestEnterpriseMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					resourceTestEnterpriseMember("deploygate_enterprise_member.current"),
				),
			},
		},
	})
}

func resourceTestEnterpriseMember(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}

		if rs.Primary.Attributes["enterprise"] == "" {
			return fmt.Errorf("enterprise expected to not be nil")
		}

		return nil
	}
}

const resourceTestEnterpriseMemberConfig = `
provider "deploygate" {
	api_key = var.user_api_key
}

variable "user_api_key" {
  type = string
}

resource "deploygate_enterprise_member" "current" {
  enterprise = var.enterprise
  members {
    name = var.add_member_name
  }
}

variable "enterprise" {
  type = string
}

variable "add_member_name" {
  type = string
}
`
