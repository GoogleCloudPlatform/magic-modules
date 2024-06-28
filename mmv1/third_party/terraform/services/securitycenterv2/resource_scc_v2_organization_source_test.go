package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSCCOrganizationSource_complete(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	suffix := acctest.RandString(t, 10)
	canonicalName := fmt.Sprintf("organizations/%s/sources/source-%s", orgId, suffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSCCOrganizationSourceCompleteExample(orgId, suffix, "My description", canonicalName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "display_name", fmt.Sprintf("TFSrc %s", suffix)),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "organization", orgId),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "description", "My description"),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "canonical_name", canonicalName),
				),
			},
			Config: testAccSCCOrganizationSourceCompleteExample(orgId, suffix, "My updated description", canonicalName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "description", "My updated description"),
				),
			},

			{
				Config: testAccSCCOrganizationSourceCompleteExample(orgId, suffix, "My updated description", canonicalName+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "canonical_name", canonicalName+"-updated"),
				),
			},

			{
				Config: testAccSCCOrganizationSourceCompleteExample(orgId, suffix+"-updated", "My updated description", canonicalName+"-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "display_name", fmt.Sprintf("TFSrc %s-updated", suffix)),
				),
			},

			{
			{
				ResourceName:      "google_scc_v2_organization_source.custom_source",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSCCOrganizationSourceCompleteExample(orgId, suffix, description, canonicalName string) string {
	return fmt.Sprintf(`
resource "google_scc_v2_organization_source" "custom_source" {
  display_name  = "TFSrc %s"
  organization  = "%s"
  description   = "%s"
  canonical_name = "%s"
}
`, suffix, orgId, description, canonicalName)
}
