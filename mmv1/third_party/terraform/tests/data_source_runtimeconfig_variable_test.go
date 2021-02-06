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
				Config: testAccRuntimeconfigDatasourceConfig(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_runtimeconfig_variable.default", "google_runtimeconfig_variable.default"),
				),
			},
		},
	})
}

func testAccRuntimeconfigDatasourceConfig(suffix string) string {
	return fmt.Sprintf(`
	resource "google_runtimeconfig_config" "default" {
		name        = "runtime-%s"
		description = "runtime-%s"
	}

	resource "google_runtimeconfig_variable" "default" {
		parent = google_runtimeconfig_config.default.name
		name   = "%s"
		value  = "%s"
	}

	data "google_runtimeconfig_variable" "default" {
		name   = google_runtimeconfig_variable.default.name
		parent = google_runtimeconfig_config.default.name
	}
`, suffix, suffix, suffix, suffix)
}
