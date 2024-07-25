package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSCCV2OrganizationSourceIamBinding(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	suffix := acctest.RandString(t, 10)
	sourceId := fmt.Sprintf("source-%s", suffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSCCV2OrganizationSourceIamBindingExample(orgId, sourceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_binding.custom_binding", "organization", orgId),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_binding.custom_binding", "source", fmt.Sprintf("organizations/%s/sources/%s", orgId, sourceId)),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source_iam_binding.custom_binding", "role", "roles/editor"),
				),
			},
			{
				ResourceName:      "google_scc_v2_organization_source_iam_binding.custom_binding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSCCV2OrganizationSourceIamBindingExample(orgId, sourceId string) string {
	return fmt.Sprintf(`
resource "google_scc_v2_organization_source" "custom_source" {
  display_name  = "TFSrc %s"
  organization  = "%s"
  description   = "Test description"
  canonical_name = "organizations/%s/sources/%s"
}

resource "google_scc_v2_organization_source_iam_binding" "custom_binding" {
  depends_on = [google_scc_v2_organization_source.custom_source]
  organization = "%s"
  source       = google_scc_v2_organization_source.custom_source.canonical_name
  role         = "roles/editor"

  members = [
    "user:example@example.com",
  ]

  condition {
    title       = "Test condition"
    description = "Test description"
    expression  = "request.time < timestamp('2025-01-01T00:00:00Z')"
  }
}
`, sourceId, orgId, orgId, sourceId, orgId)
}
