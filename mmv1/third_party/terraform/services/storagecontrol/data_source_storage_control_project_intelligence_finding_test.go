package storagecontrol_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageControlProjectIntelligenceFinding_notFound(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceGoogleStorageControlProjectIntelligenceFinding_notFound(),
				ExpectError: regexp.MustCompile(".*not found.*"),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlProjectIntelligenceFinding_notFound() string {
	return `
data "google_storage_control_project_intelligence_finding" "not_found" {
  finding_id = "nonexistent-finding-id"
}
`
}
