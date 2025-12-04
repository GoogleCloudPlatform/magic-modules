package auditmanager_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"google3/third_party/golang/terraform_providers/google_private/google_private/tpgresource/tpgresource"
)

func TestAccAuditManagerFrameworkAudit_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":               "your-org-id", // Replace with your test organization ID
		"framework_audit_id":   "my-framework-audit-update-test",
		"compliance_framework": "my-compliance-framework",
		"scope":                "initial-scope",
		"bucket_uri":           "gs://your-test-bucket", // Replace with your test bucket URI
		"updated_scope":        "updated-scope",
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuditManagerFrameworkAuditDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// First step: create the resource with the initial scope.
				Config: testAccAuditManagerFrameworkAudit_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_audit_manager_framework_audit.primary", "scope", context["scope"].(string)),
				),
			},
			{
				// Second step: update the scope and check if it was updated.
				Config: testAccAuditManagerFrameworkAudit_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_audit_manager_framework_audit.primary", "scope", context["updated_scope"].(string)),
				),
			},
		},
	})
}

// Helper function to generate the initial resource configuration.
func testAccAuditManagerFrameworkAudit_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_audit_manager_framework_audit" "primary" {
  parent               = "organizations/%{org_id}"
  framework_audit_id   = "%{framework_audit_id}"
  compliance_framework = "%{compliance_framework}"
  scope                = "%{scope}"

  framework_audit_destination {
    bucket {
      bucket_uri = "%{bucket_uri}"
    }
  }
}
`, context)
}

// Helper function to generate the updated resource configuration.
func testAccAuditManagerFrameworkAudit_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_audit_manager_framework_audit" "primary" {
  parent               = "organizations/%{org_id}"
  framework_audit_id   = "%{framework_audit_id}"
  compliance_framework = "%{compliance_framework}"
  scope                = "%{updated_scope}" // The updated value for the scope.

  framework_audit_destination {
    bucket {
      bucket_uri = "%{bucket_uri}"
    }
  }
}
`, context)
}

// You will also need a destroy check function like this one.
func testAccCheckAuditManagerFrameworkAuditDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_audit_manager_framework_audit" {
				continue
			}

			config := googleProviderConfig(t)
			url, err := tpgresource.ReplaceVars(rs.Primary, config, "{{AuditManagerBasePath}}%s")
			if err != nil {
				return err
			}
			url = fmt.Sprintf(url, rs.Primary.ID)

			_, err = transport_tpg.SendRequest(transport_tpg.TemplatedRequest{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("AuditManagerFrameworkAudit still exists at %s", url)
			}
		}

		return nil
	}
}
