package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageProjectManagementHub_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageProjectManagementHub_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_storage_project_management_hub.project_management_hub", "google_storage_project_management_hub.project_management_hub"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageProjectManagementHub_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_project_management_hub" "project_management_hub" {
  name = "%{project}"
  edition_config = "STANDARD"
}

data "google_storage_project_management_hub" "project_management_hub" {
  name = google_storage_project_management_hub.project_management_hub.name
}
`, context)
}
