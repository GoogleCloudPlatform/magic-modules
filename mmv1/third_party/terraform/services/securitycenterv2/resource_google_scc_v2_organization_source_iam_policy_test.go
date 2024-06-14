package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSCCV2OrganizationSourceIAMPolicy(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	suffix := acctest.RandString(t, 10)
	canonicalName := fmt.Sprintf("organizations/%s/sources/source-%s", orgId, suffix)
	policyData := `{"bindings":[{"role":"roles/editor","members":["user:test@example.com"]}]}`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSCCV2OrganizationSourceIAMPolicy(orgId, suffix, policyData),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_policy.custom_policy", "policy_data", policyData),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_policy.custom_policy", "source", canonicalName),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_policy.custom_policy", "organization", orgId),
				),
			},
			{
				ResourceName:      "google_scc_v2_organization_source_iam_policy.custom_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSCCV2OrganizationSourceIAMPolicy(orgId, suffix, policyData string) string {
	return fmt.Sprintf(`
resource "google_scc_v2_organization_source" "custom_source" {
  display_name  = "TFSrc %s"
  organization  = "%s"
  canonical_name = "organizations/%s/sources/source-%s"
}

resource "google_scc_v2_organization_source_iam_policy" "custom_policy" {
  organization  = "%s"
  source        = google_scc_v2_organization_source.custom_source.canonical_name
  policy_data   = "%s"
}
`, suffix, orgId, orgId, suffix, orgId, policyData)
}
