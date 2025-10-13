package cloudsecuritycompliance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCloudSecurityComplianceCloudControl_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceCloudControlConfig(t, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloudcontrol.test", "location", "global"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloudcontrol.test", "display_name", "Test CloudControl"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloudcontrol.test", "description", "A test cloud control for security compliance"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloudcontrol.test", "cloud_control_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloudcontrol.test", "name"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloudcontrol.test", "major_revision_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloudcontrol.test", "create_time"),
				),
			},
		},
	})
}

func testAccCloudSecurityComplianceCloudControlConfig(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
resource "google_cloud_security_compliance_cloudcontrol" "test" {
	organization      = "123456789"
	location          = "global"
	cloud_control_id  = "tf-test-%s"
	display_name      = "Test CloudControl"
	description       = "A test cloud control for security compliance"
	categories        = ["SECURITY"]
	severity          = "HIGH"
	finding_category  = "SECURITY_POLICY"
	remediation_steps = "Review and update the security configuration according to best practices."
	
	supported_cloud_providers        = ["GCP"]
	supported_target_resource_types = ["compute.googleapis.com/Instance"]
	
	rules {
		description         = "Ensure compute instances have secure boot enabled"
		rule_action_types   = ["DETECTIVE"]
		
		cel_expression {
			expression = "resource.data.shieldedInstanceConfig.enableSecureBoot == true"
			resource_types_values {
				values = ["compute.googleapis.com/Instance"]
			}
		}
	}
	
	parameter_spec {
		name         = "location"
		display_name = "Resource Location"
		description  = "The location where the resource should be deployed"
		value_type   = "STRING"
		is_required  = true
		
		default_value {
			string_value = "us-central1"
		}
		
		validation {
			regexp_pattern {
				pattern = "^[a-z]+-[a-z]+[0-9]$"
			}
		}
	}
}
`, suffix)
}