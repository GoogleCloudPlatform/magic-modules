package observability_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccObservabilityOrganizationSettings_datasource(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)

	context := map[string]interface{}{
		"org_id":   orgId,
		"location": "us",
	}
	dataResourceName := "data.google_observability_organization_settings.settings"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityOrganizationSettings_datasource(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataResourceName, "name"),
					resource.TestCheckResourceAttrSet(dataResourceName, "service_account_id"),
					resource.TestCheckResourceAttr(dataResourceName, "organization", orgId),
				),
			},
		},
	})
}

func testAccObservabilityOrganizationSettings_datasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_observability_organization_settings" "settings" {
	organization = "%{org_id}"
	location = "%{location}"
}
`, context)
}
