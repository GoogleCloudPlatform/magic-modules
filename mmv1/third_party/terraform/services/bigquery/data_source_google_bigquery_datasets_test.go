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
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
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
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "bigquery" {
  project = google_project.project.project_id
  service = "bigquery.googleapis.com"
}

resource "google_bigquery_dataset" "foo" {
  project 				      = google_project.project.project_id
  dataset_id                  = "tf_test_foo_%{random_suffix}"
  friendly_name               = "Foo"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000

  depends_on = [
	google_project_service.bigquery,
  ]
}

resource "google_bigquery_dataset" "bar" {
  project 				      = google_project.project.project_id
  dataset_id                  = "tf_test_bar_%{random_suffix}"
  friendly_name               = "BaR"
  description                 = "This is a test description"
  location                    = "EU"
  default_table_expiration_ms = 3600000

  depends_on = [
	google_project_service.bigquery,
  ]
}

resource "google_bigquery_dataset" "baz" {
  project 				      = google_project.project.project_id
  dataset_id                  = "tf_test_baz_%{random_suffix}"
  friendly_name               = "BaZ"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000

  depends_on = [
	google_project_service.bigquery,
  ]
}

data "google_bigquery_datasets" "example" {
  project = google_project.project.project_id
}
`, context)
}
