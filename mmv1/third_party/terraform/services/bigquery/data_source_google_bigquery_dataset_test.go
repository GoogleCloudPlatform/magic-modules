package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigquery"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigqueryreservation"
)

func TestAccDataSourceGoogleBigqueryDataset_basic(t *testing.T) {
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
				Config: testAccDataSourceGoogleBigqueryDataset_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_bigquery_dataset.bar", "google_bigquery_dataset.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBigqueryDataset_replicasAndFailoverReservation(t *testing.T) {
	t.Parallel()

	pid := envvar.GetTestProjectFromEnv()
	suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix": suffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBigqueryDataset_replicasAndFailoverReservation(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_bigquery_dataset.bar", "google_bigquery_dataset.foo"),
					resource.TestCheckResourceAttr("data.google_bigquery_dataset.bar", "failover_reservation_id", fmt.Sprintf("%s.tf-test-ds-dr-%s", pid, suffix)),
					resource.TestCheckResourceAttr("data.google_bigquery_dataset.bar", "replicas.#", "2"),
					testAccCheckBigQueryDatasetReplicasMetadata("data.google_bigquery_dataset.bar", 2, true),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBigqueryDataset_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

  resource "google_bigquery_dataset" "foo" {
    dataset_id                  = "tf_test_ds_%{random_suffix}"
    friendly_name               = "testing"
    description                 = "This is a test description"
    location                    = "US"
    default_table_expiration_ms = 3600000
  }

  data "google_bigquery_dataset" "bar" {
    dataset_id    = google_bigquery_dataset.foo.dataset_id
  }
`, context)
}

func testAccDataSourceGoogleBigqueryDataset_replicasAndFailoverReservation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "dr_res" {
  name               = "tf-test-ds-dr-%{random_suffix}"
  location           = "us-west2"
  secondary_location = "us-west1"
  slot_capacity      = 0
  edition            = "ENTERPRISE_PLUS"
  ignore_idle_slots  = true
  concurrency        = 0
  autoscale {
    max_slots = 100
  }
}

resource "google_bigquery_dataset" "foo" {
  dataset_id              = "tf_test_ds_crr_%{random_suffix}"
  friendly_name           = "testing-crr"
  description             = "CRR and DR datasource test"
  location                = "us-west2"
  failover_reservation_id = "${google_bigquery_reservation.dr_res.project}.${google_bigquery_reservation.dr_res.name}"

  replicas {
    id       = "us-west2"
    location = "us-west2"
  }

  replicas {
    id       = "us-west1"
    location = "us-west1"
  }
}

data "google_bigquery_dataset" "bar" {
  dataset_id = google_bigquery_dataset.foo.dataset_id
}
`, context)
}
