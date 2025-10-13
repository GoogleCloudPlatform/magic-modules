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
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "location", "global"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "display_name", "TF Test CloudControl"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "description", "A test cloud control for security compliance"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloud_control.test", "cloud_control_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloud_control.test", "name"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloud_control.test", "major_revision_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_cloud_control.test", "create_time"),
				),
			},
		},
	})
}

func TestAccCloudSecurityComplianceCloudControl_update(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceCloudControlConfig(t, suffix),
			},
			{
				ResourceName:      "google_cloud_security_compliance_cloud_control.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudSecurityComplianceCloudControlUpdate(t, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "display_name", "Updated TF Test CloudControl"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "description", "An updated test cloud control description for security compliance created from terraform"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "categories.#", "2"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "severity", "CRITICAL"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "finding_category", "UPDATED_SECURITY_POLICY"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "supported_cloud_providers.#", "2"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "supported_target_resource_types.#", "1"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_cloud_control.test", "parameter_spec.#", "2"),
				),
			},
			{
				ResourceName:      "google_cloud_security_compliance_cloud_control.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudSecurityComplianceCloudControlConfig(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_cloud_control" "test" {
	organization      = "1035865795181"
	name              = "organizations/1035865795181/locations/global/cloudControls/tf-test-%s"
	location          = "global"
	cloud_control_id  = "tf-test-%s"
	display_name      = "TF Test CloudControl"
	description       = "A test cloud control for security compliance"
	categories        = ["CC_CATEGORY_INFRASTRUCTURE"]
	severity          = "HIGH"
	finding_category  = "SECURITY_POLICY"
	remediation_steps = "Review and update the security configuration according to best practices."
	
	supported_cloud_providers        = ["GCP"]
	supported_target_resource_types = []
	
	rules {
		description         = "Ensure compute instances have secure boot enabled"
		rule_action_types   = ["RULE_ACTION_TYPE_DETECTIVE"]
		
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
`, suffix, suffix)
}

func testAccCloudSecurityComplianceCloudControlUpdate(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_cloud_control" "test" {
	organization      = "1035865795181"
	name              = "organizations/1035865795181/locations/global/cloudControls/tf-test-%s"
	location          = "global"
	cloud_control_id  = "tf-test-%s"
	display_name      = "Updated TF Test CloudControl"
	description       = "An updated test cloud control description for security compliance created from terraform"
	categories        = ["CC_CATEGORY_INFRASTRUCTURE", "CC_CATEGORY_SECURITY"]
	severity          = "CRITICAL"
	finding_category  = "UPDATED_SECURITY_POLICY"
	remediation_steps = "Updated remediation steps: Review and update the security configuration according to updated best practices."
	
	supported_cloud_providers        = ["GCP", "AWS"]
	supported_target_resource_types = []
	
	rules {
		description         = "Updated rule: Ensure compute instances have secure boot and integrity monitoring enabled"
		rule_action_types   = ["RULE_ACTION_TYPE_DETECTIVE", "RULE_ACTION_TYPE_PREVENTIVE"]
		
		cel_expression {
			expression = "resource.data.shieldedInstanceConfig.enableSecureBoot == true && resource.data.shieldedInstanceConfig.enableIntegrityMonitoring == true"
			resource_types_values {
				values = ["compute.googleapis.com/Instance"]
			}
		}
	}
	
	parameter_spec {
		name         = "location"
		display_name = "Updated Resource Location"
		description  = "The updated location where the resource should be deployed"
		value_type   = "STRING"
		is_required  = true
		
		default_value {
			string_value = "us-west1"
		}
		
		validation {
			regexp_pattern {
				pattern = "^[a-z]+-[a-z]+[0-9]$"
			}
		}
	}
}
`, suffix, suffix)
}