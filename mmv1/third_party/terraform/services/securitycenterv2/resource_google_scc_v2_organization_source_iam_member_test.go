package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSCCV2OrganizationSourceIAMMember(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	suffix := acctest.RandString(t, 10)
	role := "roles/editor"
	member := "user:test@example.com"
	conditionTitle := "Title"
	conditionDescription := "Description"
	conditionExpression := `request.time < timestamp(\"2023-12-31T00:00:00Z\")`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSCCV2OrganizationSourceIAMMember(orgId, suffix, role, member, conditionTitle, conditionDescription, conditionExpression),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_member.custom_member", "organization", orgId),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_member.custom_member", "role", role),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_member.custom_member", "member", member),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_member.custom_member", "condition.0.title", conditionTitle),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_member.custom_member", "condition.0.description", conditionDescription),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_member.custom_member", "condition.0.expression", conditionExpression),
				),
			},
			{
				ResourceName:      "google_scc_v2_organization_source_iam_member.custom_member",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSCCV2OrganizationSourceIAMMember(orgId, suffix, role, member, title, description, expression string) string {
	return fmt.Sprintf(`
resource "google_scc_v2_organization_source" "custom_source" {
  display_name  = "TFSrc %s"
  organization  = "%s"
  canonical_name = "organizations/%s/sources/source-%s"
}
resource "google_scc_v2_organization_source_iam_member" "custom_member" {
  organization  = "%s"
  source        = google_scc_v2_organization_source.custom_source.canonical_name
  role          = "%s"
  member        = "%s"
  condition {
    title       = "%s"
    description = "%s"
    expression  = "%s"
  }
}
`, suffix, orgId, orgId, suffix, orgId, role, member, title, description, expression)
}
