package bigqueryanalyticshub_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigquery"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigqueryanalyticshub"
)

// Dummy usages to ensure imports are satisfied in both versions
var _ = envvar.GetTestProjectFromEnv
var _ = fmt.Sprintf
var _ = os.Getenv

func TestAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubListingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListing_basic(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_update(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "log_linked_dataset_query_user_email", "true"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.restrict_direct_table_access", "true"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.restrict_query_result", "false"),
				),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_updateRestrictedExportPolicy(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.restrict_direct_table_access", "false"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.restrict_query_result", "true"),
				),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_storedProcedureConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "stored_procedure_config.0.enabled", "true"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "stored_procedure_config.0.allowed_stored_procedure_types.#", "1"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "stored_procedure_config.0.allowed_stored_procedure_types.0", "SQL_PROCEDURE"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.enabled", "false"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.restrict_direct_table_access", "false"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.restricted_export_policy.0.restrict_query_result", "false"),
				),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_storedProcedureConfigDisabled(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "stored_procedure_config.0.enabled", "false"),
				),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigqueryAnalyticsHubListing_multiregion(t *testing.T) {
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	t.Parallel()

	randomDatasetSuffix := acctest.RandString(t, 10)
	datasetID := fmt.Sprintf("tf_test_listing_replica_%s", randomDatasetSuffix)

	bqdataset, err := AddBigQueryDatasetReplica(t, envvar.GetTestProjectFromEnv(), datasetID, "us", "eu")
	if err != nil {
		t.Fatalf("Failed to create BigQuery dataset and add replica: %v", err)
	}

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"bqdataset":     bqdataset,
	}

	t.Cleanup(func() {
		CleanupBigQueryDatasetAndReplica(t, envvar.GetTestProjectFromEnv(), datasetID, "eu")
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubListingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListing_multiregion(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.replica_locations.#", "1"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "bigquery_dataset.0.replica_locations.0", "eu"),
				),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_listing.listing",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_exchange_id", "listing_id", "location"},
			},
		},
	})
}

func testAccBigqueryAnalyticsHubListing_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_de_%{random_suffix}"
  display_name     = "tf_test_de_%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_listing_%{random_suffix}"
  display_name     = "tf_test_listing_%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id = "tf_test_ds_%{random_suffix}"
  location   = "US"
}
`, context)
}

func testAccBigqueryAnalyticsHubListing_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_de_%{random_suffix}"
  display_name     = "tf_test_de_%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_listing_%{random_suffix}"
  display_name     = "tf_test_listing_%{random_suffix}"
  log_linked_dataset_query_user_email = true

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
    restricted_export_policy {
      enabled                      = true
      restrict_direct_table_access = true
      restrict_query_result        = false
    }
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id = "tf_test_ds_%{random_suffix}"
  location   = "US"
}
`, context)
}

func testAccBigqueryAnalyticsHubListing_updateRestrictedExportPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_de_%{random_suffix}"
  display_name     = "tf_test_de_%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_listing_%{random_suffix}"
  display_name     = "tf_test_listing_%{random_suffix}"
  log_linked_dataset_query_user_email = true

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
    restricted_export_policy {
      enabled                      = true
      restrict_direct_table_access = false
      restrict_query_result        = true
    }
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id = "tf_test_ds_%{random_suffix}"
  location   = "US"
}
`, context)
}

func testAccBigqueryAnalyticsHubListing_storedProcedureConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_de_%{random_suffix}"
  display_name     = "tf_test_de_%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_listing_%{random_suffix}"
  display_name     = "tf_test_listing_%{random_suffix}"
  log_linked_dataset_query_user_email = true

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
    restricted_export_policy {
      enabled                      = false
      restrict_direct_table_access = false
      restrict_query_result        = false
    }
  }

  stored_procedure_config {
    enabled = true
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id = "tf_test_ds_%{random_suffix}"
  location   = "US"
}
`, context)
}

func testAccBigqueryAnalyticsHubListing_storedProcedureConfigDisabled(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_de_%{random_suffix}"
  display_name     = "tf_test_de_%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_listing_%{random_suffix}"
  display_name     = "tf_test_listing_%{random_suffix}"
  log_linked_dataset_query_user_email = true

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }

  stored_procedure_config {
    enabled = false
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id = "tf_test_ds_%{random_suffix}"
  location   = "US"
}
`, context)
}

func testAccBigqueryAnalyticsHubListing_multiregion(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "us"
  data_exchange_id = "tf_test_de_%{random_suffix}"
  display_name     = "tf_test_de_%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "us"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_listing_%{random_suffix}"
  display_name     = "tf_test_listing_%{random_suffix}"

  bigquery_dataset {
    dataset = "%{bqdataset}"
    replica_locations = ["eu"]
  }
}
`, context)
}
