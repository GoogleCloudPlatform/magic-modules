package storagecontrol_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageControlFolderIntelligenceFindingsSummary_empty(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageControlFolderIntelligenceFindingsSummary_empty(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_storage_control_folder_intelligence_findings_summary.empty", "total_findings_count", "0"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageControlFolderIntelligenceFindingsSummary_empty(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent              = "organizations/%{org_id}"
  display_name        = "tf-test-folder-name%{random_suffix}"
  deletion_protection = false
}

data "google_storage_control_folder_intelligence_findings_summary" "empty" {
  folder     = google_folder.folder.folder_id
  depends_on = [google_folder.folder]
}
`, context)
}
