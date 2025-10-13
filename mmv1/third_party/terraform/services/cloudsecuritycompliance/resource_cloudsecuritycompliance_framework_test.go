package cloudsecuritycompliance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCloudSecurityComplianceFramework_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceFrameworkConfig(t, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloudsecuritycompliance_framework.test", "location", "global"),
					resource.TestCheckResourceAttr("google_cloudsecuritycompliance_framework.test", "display_name", "Test Framework"),
					resource.TestCheckResourceAttr("google_cloudsecuritycompliance_framework.test", "description", "A test framework for cloud security compliance"),
					resource.TestCheckResourceAttrSet("google_cloudsecuritycompliance_framework.test", "framework_id"),
					resource.TestCheckResourceAttrSet("google_cloudsecuritycompliance_framework.test", "name"),
					resource.TestCheckResourceAttrSet("google_cloudsecuritycompliance_framework.test", "major_revision_id"),
					resource.TestCheckResourceAttrSet("google_cloudsecuritycompliance_framework.test", "type"),
				),
			},
		},
	})
}

func testAccCloudSecurityComplianceFrameworkConfig(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
resource "google_cloud_security_compliance_framework" "test" {
	organization   = "123456789"
	location       = "global"
	framework_id   = "tf-test-%s"
	display_name   = "Test Framework"
	description    = "A test framework for cloud security compliance"
	category       = ["CC_CATEGORY_INFRASTRUCTURE"]
}
`, suffix)
}
