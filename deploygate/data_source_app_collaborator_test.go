package deploygate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_AppCollaborator_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { Test_DGPreCheck(t) },
		Providers: testDGProviders,
		Steps: []resource.TestStep{
			{
				Config: testAppCollaboratorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAppCollaborator("data.deploygate_app_collaborator.current"),
				),
			},
		},
	})
}

func testAppCollaborator(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find app users resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Account Id resource ID not set")
		}

		if rs.Primary.Attributes["users"] == "" {
			return fmt.Errorf("Users expected to not be nil")
		}

		if rs.Primary.Attributes["teams"] == "" {
			return fmt.Errorf("Teams expected to not be nil")
		}

		if rs.Primary.Attributes["usage"] == "" {
			return fmt.Errorf("Usage expected to not be nil")
		}

		return nil
	}
}

const testAppCollaboratorConfig = `
data "deploygate_app_collaborator" "current" {}
`
