package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigQueryDatasetListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetListResource_basic(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "dataset_id", datasetID),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "project", project),
				),
			},
			{
				Query:  true,
				Config: testAccBigQueryDatasetListResource_query(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_bigquery_dataset.all", 1),
					querycheck.ExpectIdentity("google_bigquery_dataset.all", map[string]knownvalue.Check{
						"dataset_id": knownvalue.StringExact(datasetID),
						"project":    knownvalue.StringExact(project),
					}),
				},
			},
		},
	})
}

func testAccBigQueryDatasetListResource_basic(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id = %q
  location   = "US"
}
`, datasetID)
}

func testAccBigQueryDatasetListResource_query(project string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_bigquery_dataset" "all" {
  provider = google
	limit    = 10000

  config {
    project = %q
  }
}
`, project)
}
