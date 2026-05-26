package storagecontrol_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageControlProjectIntelligenceFindingRevisions_notFound(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleStorageControlProjectIntelligenceFindingRevisions_notFound(),
				ExpectError: regexp.MustCompile(".*not found.*"),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlProjectIntelligenceFindingRevisions_notFound() string {
	return `
data "google_storage_control_project_intelligence_finding_revisions" "not_found" {
  finding_id = "nonexistent-finding-id"
}
`
}
