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
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.#", "1"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.dataset_id", fmt.Sprintf("tf_test_foo_%s", context["random_suffix"])),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.labels.%", "1"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.labels.goog-terraform-provisioned", "true"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.location", "US"),
					resource.TestCheckResourceAttr("data.google_bigquery_datasets.example", "datasets.0.friendly_name", "Foo"),
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

  access {
    role          = "OWNER"
    user_by_email = google_service_account.bqowner.email
  }
}

resource "google_service_account" "bqowner" {
  account_id = "bqowner"
}

data "google_bigquery_datasets" "example" {
  project = "%{project_id}"
  depends_on = [
	google_bigquery_dataset.foo,
  ]
}
`, context)
}
