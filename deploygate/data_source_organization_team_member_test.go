package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_DataSourceOrganizationTeamMember_basic(t *testing.T) {
	testWithVCR(t, resource.TestCase{
		PreCheck: func() { Test_DGPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationTeamMemberConfig,
				Check: resource.ComposeTestCheckFunc(
					testDataSourceOrganizationTeamMember("data.deploygate_organization_team_member.current"),
				),
			},
		},
	})
}

func testDataSourceOrganizationTeamMember(n string) resource.TestCheckFunc {
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

		return nil
	}
}

const testDataSourceOrganizationTeamMemberConfig = `
provider "deploygate" {
	api_key = "organization_api_key"
}

data "deploygate_organization_team_member" "current" {
	organization = "test-organization-fnaoto"
	team 				 = "admin"
}
`
