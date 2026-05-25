package storagecontrol_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageControlProjectIntelligenceFindingsSummary_empty(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleStorageControlProjectIntelligenceFindingsSummary_empty(),
				ExpectError: regexp.MustCompile(".*not found.*"),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlProjectIntelligenceFindingsSummary_empty() string {
	return `
data "google_storage_control_project_intelligence_findings_summary" "empty" {
}
`
}
