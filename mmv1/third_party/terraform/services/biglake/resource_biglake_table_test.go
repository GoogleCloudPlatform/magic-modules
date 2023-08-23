package biglake_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBiglakeTable_bigqueryBiglakeTable_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeTable_bigqueryBiglakeTableExample(context),
			},
			{
				ResourceName:            "google_biglake_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "database_id", "catalog_id"},
			},
			{
				Config: testAccBiglakeTable_bigqueryBiglakeTable_update(context),
			},
			{
				ResourceName:            "google_biglake_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "database_id", "catalog_id"},
			},
		},
	})
}

func testAccBiglakeTable_bigqueryBiglakeTable_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_biglake_table" "table" {
    name = "tf_test_my_table%{random_suffix}"
    catalog_id = google_biglake_catalog.catalog.catalog_id
    database_id = google_biglake_database.database.database_id
    location = google_biglake_catalog.catalog.location
    type = "HIVE"
    hive_options {
      table_type = "EXTERNAL_TABLE"
      storage_descriptor {
        location_uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.data_folder.name}/data"
      }
	  parameters = {
		"name" = "hammer"
		"drink" = "lemonade"
	  }
    }
}
`, context)
}
