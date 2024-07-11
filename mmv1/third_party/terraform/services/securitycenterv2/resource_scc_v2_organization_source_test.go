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
	displayName := fmt.Sprintf("TFSrc %s", suffix)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccTestPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSCCOrganizationSourceCompleteExample(orgId, suffix, "My description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "display_name", displayName),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "organization", orgId),
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "description", "My description"),
					resource.TestCheckResourceAttrset("google_scc_v2_organization_source.custom_source", "canonical_name"),
				),
			},
			{
				Config: testAccSCCOrganizationSourceCompleteExample(orgId, suffix, "My updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("google_scc_v2_organization_source.custom_source", "description", "My updated description"),
				),
			},
			{
				ResourceName:      "google_scc_v2_organization_source.custom_source",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSCCOrganizationSourceCompleteExample(orgId, suffix, description string) string {
	return fmt.Sprintf(`
resource "google_scc_v2_organization_source" "custom_source" {
  display_name  = "TFSrc %s"
  organization  = "%s"
  description   = "%s"

}
output "canonical_name" {
  value = google_scc_v2_organization_source.custom_source.canonical_name
}  
`, suffix, orgId, description)
}
