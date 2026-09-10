package datacatalog_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/datacatalog"
)

func TestAccDataSourceGoogleDataCatalogTaxonomy_basic(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleDataCatalogTaxonomy_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.google_data_catalog_taxonomy.test", "name", "google_data_catalog_taxonomy.test", "name"),
					resource.TestCheckResourceAttrPair("data.google_data_catalog_taxonomy.test", "description", "google_data_catalog_taxonomy.test", "description"),
					resource.TestCheckResourceAttrPair("data.google_data_catalog_taxonomy.test", "activated_policy_types.#", "google_data_catalog_taxonomy.test", "activated_policy_types.#"),
					resource.TestCheckResourceAttrPair("data.google_data_catalog_taxonomy.test", "project", "google_data_catalog_taxonomy.test", "project"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleDataCatalogTaxonomy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_data_catalog_taxonomy" "test" {
  display_name           = "tf_test_%{random_suffix}"
  description            = "A test taxonomy"
  activated_policy_types = ["FINE_GRAINED_ACCESS_CONTROL"]
  region                 = "us-central1"
}

data "google_data_catalog_taxonomy" "test" {
  display_name = google_data_catalog_taxonomy.test.display_name
  region       = google_data_catalog_taxonomy.test.region
}
`, context)
}
