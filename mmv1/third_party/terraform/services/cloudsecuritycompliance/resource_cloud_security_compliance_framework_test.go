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
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "location", "global"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "display_name", "Test Framework"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "description", "A test framework for cloud security compliance"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework.test", "framework_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework.test", "name"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework.test", "major_revision_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework.test", "type"),
				),
			},
		},
	})
}

func TestAccCloudSecurityComplianceFramework_update(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceFrameworkConfig(t, suffix),
			},
			{
				ResourceName:      "google_cloud_security_compliance_framework.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudSecurityComplianceFrameworkUpdate(t, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "display_name", "Updated TF Test Framework"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "description", "An updated test framework description for cloud security compliance created from terraform"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "category.#", "2"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework.test", "cloud_control_details.#", "2"),
				),
			},
			{
				ResourceName:      "google_cloud_security_compliance_framework.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudSecurityComplianceFrameworkConfig(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_framework" "test" {
	organization   = "1035865795181"
	name 		   = "organizations/1035865795181/locations/global/frameworks/tf-test-%s"
	location       = "global"
	framework_id   = "tf-test-%s"
	display_name   = "Created TF Test Framework"
	description    = "A test framework for cloud security compliance created from terraform"
	category       = ["CUSTOM_FRAMEWORK"]
	
	cloud_control_details {
		name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-resource-availability"
		major_revision_id = "1"
	}
}
`, suffix, suffix)
}

func testAccCloudSecurityComplianceFrameworkUpdate(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_framework" "test" {
	organization   = "1035865795181"
	name 		   = "organizations/1035865795181/locations/global/frameworks/tf-test-%s"
	location       = "global"
	framework_id   = "tf-test-%s"
	display_name   = "Updated TF Test Framework"
	description    = "An updated test framework description for cloud security compliance created from terraform"
	category       = ["CUSTOM_FRAMEWORK"]
}
`, suffix, suffix)
}
