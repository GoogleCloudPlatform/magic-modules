package storagecontrol_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageControlFolderIntelligenceFindingsSummary_empty(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id": envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageControlFolderIntelligenceFindingsSummary_empty(context),
				ExpectError: regexp.MustCompile(".*not found.*"),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlFolderIntelligenceFindingsSummary_empty(context map[string]interface{}) string {
	return acctest.Nprintf(`

data "google_storage_control_folder_intelligence_findings_summary" "empty" {
  folder = "folders/%{org_id}"
}
`, context)
}
