package storagecontrol_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageControlProjectIntelligenceFindingRevision_notFound(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleStorageControlProjectIntelligenceFindingRevision_notFound(),
				ExpectError: regexp.MustCompile(".*not found.*"),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlProjectIntelligenceFindingRevision_notFound() string {
	return `
data "google_storage_control_project_intelligence_finding_revision" "not_found" {
  finding_id  = "nonexistent-finding-id"
  revision_id = "nonexistent-revision-id"
}
`
}
