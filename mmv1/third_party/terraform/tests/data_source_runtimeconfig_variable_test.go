package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRuntimeconfigVariableDatasource_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRuntimeconfigDatasourceVariable(randString(t, 10), randString(t, 10), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_runtimeconfig_variable.default", "google_runtimeconfig_variable.default"),
				),
			},
		},
	})
}

func testAccRuntimeconfigDatasourceVariable(suffix string, name string, text string) string {
	return fmt.Sprintf(`
	resource "google_project_service" "default" {
		service = "runtimeconfig.googleapis.com"
	}

	resource "google_runtimeconfig_config" "default" {
		project     = google_project_service.default.project
		name        = "runtime-%s"
		description = "runtime-%s"
		depends_on  = [ google_project_service.default ]
	}

	resource "google_runtimeconfig_variable" "default" {
		project = google_project_service.default.project
		parent  = google_runtimeconfig_config.default.name
		name    = "%s"
		text    = "%s"
	}

	data "google_runtimeconfig_variable" "default" {
		project = google_project_service.default.project
		name    = google_runtimeconfig_variable.default.name
		parent  = google_runtimeconfig_config.default.name
	}
`, suffix, suffix, name, text)
}
