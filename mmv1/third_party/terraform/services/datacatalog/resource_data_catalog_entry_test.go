package datacatalog_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataCatalogEntry_update(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":        envvar.GetTestProjectFromEnv(),
		"location":       "us-central1",
		"random_suffix":  randString,
		"entry_id":       "tf_test_my_entry" + randString,
		"entry_group_id": "tf_test_my_group" + randString,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataCatalogEntryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCatalogEntry_dataCatalogEntryBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry.basic_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntry_dataCatalogEntryFullExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry.basic_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataCatalogEntry_dataCatalogEntryBasicExample(context),
			},
			{
				ResourceName:      "google_data_catalog_entry.basic_entry",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
