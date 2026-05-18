package licensemanager_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccLicenseManagerConfiguration_lifecycle(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"configuration_id": "tf-test-example-config-" + randomSuffix,
		"product":          "Office2021ProfessionalPlus",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseManagerConfiguration_active(context, true, 10),
			},
			{
				Config: testAccLicenseManagerConfiguration_active(context, true, 15),
			},
			{
				Config: testAccLicenseManagerConfiguration_active(context, false, 15),
			},
			{
				Config: testAccLicenseManagerConfiguration_active(context, true, 15),
			},
		},
	})
}

func testAccLicenseManagerConfiguration_active(context map[string]interface{}, active bool, licensecount int) string {
	return acctest.Nprintf(`
resource "google_license_manager_configuration" "example" {
  location         = "us-central1"
  configuration_id = "%{configuration_id}"
  product          = "%{product}"
  licensecount     = `+fmt.Sprintf("%d", licensecount)+`
  active           = `+fmt.Sprintf("%t", active)+`
}
`, context)
}
