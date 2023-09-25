package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterOrganizationCustomModule_sccOrganizationCustomModuleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterOrganizationCustomModuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterOrganizationCustomModule_sccOrganizationCustomModuleFullExample(context),
			},
			{
				ResourceName:      "google_scc_organization_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityCenterOrganizationCustomModule_sccOrganizationCustomModuleUpdate(context),
			},
			{
				ResourceName:      "google_scc_organization_custom_module.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSecurityCenterOrganizationCustomModule_sccOrganizationCustomModuleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_organization_custom_module" "example" {
	organization = "%{org_id}"
	display_name = "tf_test_full_custom_module%{random_suffix}"
	enablement_state = "DISABLED"
	custom_config {
		predicate {
			expression = "resource.name == \"updated-name\""
			title = "Updated expression title"
			description = "Updated description of the expression"
			location = "Updated location of the expression"
		}
		custom_output {
			properties {
				name = "violation"
				value_expression {
					expression = "resource.name"
					title = "Updated expression title"
					description = "Updated description of the expression"
					location = "Updated location of the expression"
				}
			}
		}
		resource_selector {
			resource_types = [
				"compute.googleapis.com/Instance",
			]
		}
		severity = "CRITICAL"
		description = "Updated description of the custom module"
		recommendation = "Updated steps to resolve violation"
	}
}
`, context)
}
