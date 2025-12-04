package auditmanager_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
)

func TestAccAuditManagerFrameworkAuditScopeReport_basic(t *testing.T) {
	t.Parallel()
	s
	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditManagerFrameworkAuditScopeReport_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_audit_manager_framework_audit_scope_report.report", "name", "organizations/"+context["org_id"].(string)+"/locations/global/frameworkAuditScopeReports/my-compliance-framework"),
					resource.TestCheckResourceAttr("data.google_audit_manager_framework_audit_scope_report.report", "compliance_framework", "my-compliance-framework"),
				),
			},
		},
	})
}

func testAccAuditManagerFrameworkAuditScopeReport_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_audit_manager_framework_audit_scope_report" "report" {
  scope                = "organizations/%{org_id}/locations/global"
  report_format        = "ODF"
  compliance_framework = "my-compliance-framework"
}
`, context)
}
