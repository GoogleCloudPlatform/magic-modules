package cloudsecuritycompliance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCloudSecurityComplianceCloudControl_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceCloudControl_basic(context),
			},
			{
				ResourceName:            "google_cloud_security_compliance_cloud_control.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_control_id", "location", "organization"},
			},
		},
	})
}

func testAccCloudSecurityComplianceCloudControl_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_cloud_control" "test" {
	organization      = "%{org_id}"
	name              = "organizations/%{org_id}/locations/global/cloudControls/tf-test-%{random_suffix}"
	location          = "global"
	cloud_control_id  = "tf-test-%{random_suffix}"
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
`, context)
}

func TestAccCloudSecurityComplianceCloudControl_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceCloudControl_basic(context),
			},
			{
				ResourceName:            "google_cloud_security_compliance_cloud_control.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_control_id", "location", "organization"},
			},
			{
				Config: testAccCloudSecurityComplianceCloudControl_update(context),
			},
			{
				ResourceName:            "google_cloud_security_compliance_cloud_control.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_control_id", "location", "organization"},
			},
		},
	})
}

func testAccCloudSecurityComplianceCloudControl_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_cloud_control" "example" {
  name              = "organizations/%{org_id}/locations/global/cloudControls/tf-test-%{random_suffix}"
  organization      = "%{org_id}"
  location          = "global"
  cloud_control_id  = "tf-test-%{random_suffix}"

  
  display_name      = "Updated CloudControl Name"
  description       = "An updated description for the cloud control"
  categories        = ["CC_CATEGORY_INFRASTRUCTURE"]
  severity          = "CRITICAL"
  finding_category  = "UPDATED_SECURITY_POLICY"
  remediation_steps = "Updated remediation steps with more detailed instructions for security configuration."
  
  supported_cloud_providers        = ["GCP", "AWS"]
  supported_target_resource_types = ["TARGET_RESOURCE_CRM_TYPE_ORG"]
  
  rules {
    description         = "Updated rule: Ensure compute instances have secure boot and integrity monitoring enabled"
    rule_action_types   = ["RULE_ACTION_TYPE_DETECTIVE"]
    
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
    description  = "Updated description for the location parameter"
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
`, context)
}
