package biglake_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBiglakeDatabase_bigqueryBiglakeDatabase_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBiglakeDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBiglakeDatabase_bigqueryBiglakeDatabaseExample(context),
			},
			{
				ResourceName:            "google_biglake_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "catalog_id"},
			},
			{
				Config: testAccBiglakeDatabase_bigqueryBiglakeDatabase_update(context),
			},
			{
				ResourceName:            "google_biglake_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "catalog_id"},
			},
		},
	})
}

func testAccBiglakeDatabase_bigqueryBiglakeDatabase_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_biglake_database" "database" {
	# Update Database Id
    name = "tf-test-my-database%{random_suffix}"
    catalog_id = google_biglake_catalog.default.name
    # Hard code to avoid invalid random id suffix
    location = google_biglake_catalog.catalog.location
	type = "HIVE"
	hive_options {
        location_uri = "gs://${google_storage_bucket.bucket.name}/${google_storage_bucket_object.metadata_folder.name}/metadata"
		parameters = {
			"tool" = "screwdriver"
		}
    }
}
`, context)
}
