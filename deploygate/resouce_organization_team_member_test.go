package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_ResourceOrganizationTeamMember_basic(t *testing.T) {
	testWithVCR(t, resource.TestCase{
		PreCheck: func() { Test_DGPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: resourceTestOrganizationTeamMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					resourceTestOrganizationTeamMember("deploygate_organization_team_member.current"),
				),
			},
		},
	})
}

func resourceTestOrganizationTeamMember(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}

		if rs.Primary.Attributes["organization"] == "" {
			return fmt.Errorf("organization expected to not be nil")
		}

		if rs.Primary.Attributes["team"] == "" {
			return fmt.Errorf("team expected to not be nil")
		}

		if rs.Primary.Attributes["users"] == "" {
			return fmt.Errorf("users expected to not be nil")
		}

		return nil
	}
}

const resourceTestOrganizationTeamMemberConfig = `
provider "deploygate" {
	api_key = var.organization_api_key
}

variable "organization_api_key" {
  type = string
}

resource "deploygate_organization_team_member" "current" {
  organization = var.organization
	team 				 = var.team

  users {
    name = var.add_member_name
  }
}

variable "organization" {
  type = string
}

variable "team" {
  type = string
}

variable "add_member_name" {
  type = string
}
`
