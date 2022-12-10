package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_ResourceEnterpriseOrganizationMember_basic(t *testing.T) {
	testWithVCR(t, resource.TestCase{
		PreCheck: func() { Test_DGPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: resourceTestEnterpriseOrganizationMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					resourceTestEnterpriseOrganizationMember("deploygate_enterprise_organization_member.current"),
				),
			},
		},
	})
}

func resourceTestEnterpriseOrganizationMember(n string) resource.TestCheckFunc {
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

		if rs.Primary.Attributes["organization"] == "" {
			return fmt.Errorf("organization expected to not be nil")
		}

		return nil
	}
}

const resourceTestEnterpriseOrganizationMemberConfig = `
provider "deploygate" {
	api_key = "enterprise_api_key"
}

resource "deploygate_enterprise_organization_member" "current" {
  enterprise   = "test-enterprise-fnaoto"
	organization = "test-organization-fnaoto"

  users {
    name = "naoto-fukuda"
  }
}
`
