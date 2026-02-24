package observability_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccObservabilityFolderSettings_datasource(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	folderDisplayName := "tf-test-" + acctest.RandString(t, 10)

	context := map[string]interface{}{
		"org_id":              orgId,
		"folder_display_name": folderDisplayName,
		"location":            "us",
	}
	dataResourceName := "data.google_observability_folder_settings.settings"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityFolderSettings_datasource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataResourceName, "name"),
					resource.TestCheckResourceAttrSet(dataResourceName, "service_account_id"),
					resource.TestCheckResourceAttrPair(dataResourceName, "folder", "google_folder.test", "folder_id"),
				),
			},
		},
	})
}

func testAccObservabilityFolderSettings_datasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "test" {
	display_name        = "%{folder_display_name}"
	parent              = "organizations/%{org_id}"
	deletion_protection = false
}

resource "time_sleep" "wait_for_folder" {
	create_duration = "90s"
	depends_on      = [google_folder.test]
}

data "google_observability_folder_settings" "settings" {
	folder     = google_folder.test.folder_id
	location   = "%{location}"
	depends_on = [time_sleep.wait_for_folder]
}
`, context)
}
