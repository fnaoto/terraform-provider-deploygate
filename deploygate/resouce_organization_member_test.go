package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_ResourceOrganizationMember_basic(t *testing.T) {
	testWithVCR(t, resource.TestCase{
		PreCheck: func() { Test_DGPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: resourceTestOrganizationMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					resourceTestOrganizationMember("deploygate_organization_member.current"),
				),
			},
		},
	})
}

func resourceTestOrganizationMember(n string) resource.TestCheckFunc {
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

		return nil
	}
}

const resourceTestOrganizationMemberConfig = `
provider "deploygate" {
	api_key = "organization_api_key"
}

resource "deploygate_organization_member" "current" {
  organization = "test-organization-fnaoto"

  members {
    name = "naoto-fukuda"
  }
}
`
