package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleBigqueryDatasets_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project_id":    envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryDatasets_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.#", "3"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.dataset_id", fmt.Sprintf("tf_test_foo_%s", context["random_suffix"])),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.labels.%", "1"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.location", "US"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.friendly_name", "Foo"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.1.dataset_id", fmt.Sprintf("tf_test_bar_%s", context["random_suffix"])),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.1.friendly_name", "BaR"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.1.location", "EU"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.2.dataset_id", fmt.Sprintf("tf_test_baz_%s", context["random_suffix"])),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.2.friendly_name", "BaZ"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryDatasets_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "foo" {
  dataset_id                  = "tf_test_foo_%{random_suffix}"
  friendly_name               = "Foo"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_bar_%{random_suffix}"
  friendly_name               = "BaR"
  description                 = "This is a test description"
  location                    = "EU"
  default_table_expiration_ms = 3600000
}

resource "google_bigquery_dataset" "baz" {
  dataset_id                  = "tf_test_baz_%{random_suffix}"
  friendly_name               = "BaZ"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000
}

data "google_bigquery_datasets" "example" {
  project = "%{project_id}"
  depends_on = [
	google_bigquery_dataset.foo,
	google_bigquery_dataset.bar,
	google_bigquery_dataset.baz,
  ]
}
`, context)
}
