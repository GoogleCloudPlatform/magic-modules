package cloudsecuritycompliance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func testAccCloudSecurityComplianceFramework_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_security_compliance_framework" "example" {
  organization = "%{org_id}"
  location     = "global"
  framework_id = "tf-test-example-framework%{random_suffix}"
  
  display_name = "Terraform Framework Name"
  description  = "An Terraform description for the framework"
  
  cloud_control_details {
		name              = "organizations/%{org_id}/locations/global/cloudControls/builtin-assess-resource-availability"
		major_revision_id = "1"
    
    parameters {
      name = "location"
      parameter_value {
        string_value = "us-central1"
      }
    }
    parameters {
      name = "oneof-parameter"
      parameter_value {
        oneof_value {
          name = "test-oneof"
          parameter_value {
            string_value = "test-value"
          }
        }
      }
    }
  }
}
`, context)
}

func TestAccCloudSecurityComplianceFramework_update(t *testing.T) {
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
				Config: testAccCloudSecurityComplianceFramework_basic(context),
			},
			{
				ResourceName:            "google_cloud_security_compliance_framework.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"framework_id", "location", "organization"},
			},
			{
				Config: testAccCloudSecurityComplianceFramework_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_cloud_security_compliance_framework.example", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_cloud_security_compliance_framework.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"framework_id", "location", "organization"},
			},
		},
	})
}

func testAccCloudSecurityComplianceFramework_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_security_compliance_framework" "example" {
  organization = "%{org_id}"
  location     = "global"
  framework_id = "tf-test-example-framework%{random_suffix}"
  
  display_name = "Updated Terraform Framework Name"
  description  = "An updated description for the framework with additional details"
  
  cloud_control_details {
    name              = "organizations/%{org_id}/locations/global/cloudControls/builtin-data-access-governance"
    major_revision_id = "1"
    
    parameters {
      name = "region"
      parameter_value {
        string_value = "eu"
      }
    }
    parameters {
      name = "oneof-parameter"
      parameter_value {
        oneof_value {
          name = "updated-oneof"
          parameter_value {
            string_value = "updated-value"
          }
        }
      }
    }
  }
}
`, context)
}
