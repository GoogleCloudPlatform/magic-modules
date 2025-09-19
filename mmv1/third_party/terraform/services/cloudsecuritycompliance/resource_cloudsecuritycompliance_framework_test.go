package cloudsecuritycompliance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance test for Framework resource
func TestAccCloudSecurityComplianceFramework_basic(t *testing.T) {
   resourceName := "google_cloudsecuritycompliance_framework.test"
   location := "us-central1"
   name := "test-framework"

   testAccFrameworkConfig := fmt.Sprintf(`
resource "google_cloudsecuritycompliance_framework" "test" {
   location = "%s"
   name     = "%s"
}
`, location, name)

   resource.ParallelTest(t, resource.TestCase{
      PreCheck:          func() { testAccPreCheck(t) },
      ProviderFactories: testAccProviderFactories,
      Steps: []resource.TestStep{{
         Config: testAccFrameworkConfig,
         Check: resource.ComposeTestCheckFunc(
            resource.TestCheckResourceAttr(resourceName, "location", location),
            resource.TestCheckResourceAttr(resourceName, "name", name),
         ),
      }},
   })
}
