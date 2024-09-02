package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBigqueryDatasets_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryDatasets_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.#", "2"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.dataset_id", fmt.Sprintf("tf_test_ds_%s_%s", context["random_suffix"], "0")),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.1.dataset_id", fmt.Sprintf("tf_test_ds_%s_%s", context["random_suffix"], "1")),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryDatasets_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
  
  resource "google_bigquery_dataset" "test" {
    count                       = 2
    dataset_id                  = "tf_test_ds_%{random_suffix}_${count.index}"
    friendly_name               = "testing"
    description                 = "This is a test description"
    location                    = "US"
    default_table_expiration_ms = 3600000
  }

  data "google_bigquery_datasets" "example" {
    depends_on = [
      google_bigquery_dataset.test
    ]
  }
`, context)
}
