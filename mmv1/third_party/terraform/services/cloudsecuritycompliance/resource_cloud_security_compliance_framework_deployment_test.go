package cloudsecuritycompliance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCloudSecurityComplianceFrameworkDeployment_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceFrameworkDeploymentBasicConfig(t, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework_deployment.test", "location", "global"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework_deployment.test", "description", "A test framework deployment for cloud security compliance"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "framework_deployment_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "name"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "create_time"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "deployment_state"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "computed_target_resource"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework_deployment.test", "cloud_control_metadata.#", "1"),
				),
			},
			{
				ResourceName:      "google_cloud_security_compliance_framework_deployment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudSecurityComplianceFrameworkDeployment_full(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityComplianceFrameworkDeploymentFullConfig(t, acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework_deployment.test", "location", "global"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework_deployment.test", "description", "A test framework deployment with target resource creation"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "framework_deployment_id"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "name"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "create_time"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "deployment_state"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "computed_target_resource"),
					resource.TestCheckResourceAttrSet("google_cloud_security_compliance_framework_deployment.test", "target_resource_display_name"),
					resource.TestCheckResourceAttr("google_cloud_security_compliance_framework_deployment.test", "cloud_control_metadata.#", "2"),
				),
			},
			{
				ResourceName:      "google_cloud_security_compliance_framework_deployment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudSecurityComplianceFrameworkDeploymentBasicConfig(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_framework" "test_framework" {
	organization   = "1035865795181"
	name = "organizations/1035865795181/locations/global/frameworks/tf-test-framework-%s"
	location       = "global"
	framework_id   = "tf-test-framework-%s"
	display_name   = "Test Framework for Deployment"
	description    = "A test framework for deployment testing"
	category       = ["CUSTOM_FRAMEWORK"]
	
	cloud_control_details {
		name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-resource-availability"
		major_revision_id = "1"
	}
}

resource "google_cloud_security_compliance_framework_deployment" "test" {
	organization              = "1035865795181"
	location                  = "global"
	framework_deployment_id   = "tf-test-deployment-%s"
	description               = "A test framework deployment for cloud security compliance"
	
	framework {
		framework          = google_cloud_security_compliance_framework.test_framework.name
		major_revision_id  = "1"
	}
	
	target_resource_config {
		existing_target_resource = "organizations/1035865795181"
	}
	
	cloud_control_metadata {
		enforcement_mode = "DETECTIVE"
		
		cloud_control_details {
			name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-resource-availability"
			major_revision_id = "1"
			
			parameters {
				name = "location"
				parameter_value {
					string_value = "us-central1"
				}
			}
		}
	}
}
`, suffix, suffix, suffix)
}

func testAccCloudSecurityComplianceFrameworkDeploymentFullConfig(t *testing.T, suffix string) string {
	return fmt.Sprintf(`
provider "google" {
  billing_project = "terraform-test-472610"
  user_project_override = true
}

resource "google_cloud_security_compliance_framework" "test_framework" {
	organization   = "1035865795181"
	location       = "global"
	framework_id   = "tf-test-framework-%s"
	display_name   = "Test Framework for Deployment"
	description    = "A test framework for deployment testing"
	category       = ["CUSTOM_FRAMEWORK"]
	
	cloud_control_details {
		name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-resource-availability"
		major_revision_id = "1"
	}
	
	cloud_control_details {
		name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-network-security"
		major_revision_id = "1"
	}
}

resource "google_cloud_security_compliance_framework_deployment" "test" {
	organization              = "1035865795181"
	location                  = "global"
	framework_deployment_id   = "tf-test-deployment-%s"
	description               = "A test framework deployment with target resource creation"
	
	framework {
		framework          = google_cloud_security_compliance_framework.test_framework.name
		major_revision_id  = "1"
	}
	
	target_resource_config {
		target_resource_creation_config {
			folder_creation_config {
				folder_display_name = "Test Folder for Deployment %s"
				parent              = "organizations/1035865795181"
			}
		}
	}
	
	cloud_control_metadata {
		enforcement_mode = "DETECTIVE"
		
		cloud_control_details {
			name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-resource-availability"
			major_revision_id = "1"
			
			parameters {
				name = "location"
				parameter_value {
					string_value = "us-central1"
				}
			}
		}
	}
	
	cloud_control_metadata {
		enforcement_mode = "AUDIT"
		
		cloud_control_details {
			name              = "organizations/1035865795181/locations/global/cloudControls/builtin-assess-network-security"
			major_revision_id = "1"
			
			parameters {
				name = "environment"
				parameter_value {
					string_value = "production"
				}
			}
			
			parameters {
				name = "enable_monitoring"
				parameter_value {
					bool_value = true
				}
			}
		}
	}
}
`, suffix, suffix, suffix)
}