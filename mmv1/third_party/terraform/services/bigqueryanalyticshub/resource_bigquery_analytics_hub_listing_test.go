package bigqueryanalyticshub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

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
				Config: testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingBasicExample(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingUpdate(context),
				Check: resource.ComposeTestCheckFunc(
					// Verify log_linked_dataset_query_user_email has been set to true (at top level)
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "log_linked_dataset_query_user_email", "true"),
				),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListing_pubsubListingUpdateConfig(context, `["us-central1"]`, "Example for pubsub topic source - initial"),
				Check: resource.ComposeTestCheckFunc(
					// Verify initial state for Pub/Sub listing
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing_pubsub", "pubsub_topic.0.data_affinity_regions.#", "1"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing_pubsub", "pubsub_topic.0.data_affinity_regions.0", "us-central1"),
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing_pubsub", "description", "Example for pubsub topic source - initial"),
				),
			},
			// Step 7: Import the updated Pub/Sub Topic listing to verify import after update.
			{
				ResourceName:      "google_bigquery_analytics_hub_listing.listing_pubsub",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "US"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}"
  description      = "example data exchange%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example data exchange update%{random_suffix}"
  log_linked_dataset_query_user_email  = true

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}
`, context)
}
func testAccBigqueryAnalyticsHubListing_pubsubListingUpdateConfig(context map[string]interface{}, dataAffinityRegionsHCL string, description string) string {
	// Create a mutable copy of the context map
	updatedContext := make(map[string]interface{})
	for k, v := range context {
		updatedContext[k] = v
	}

	// Directly assign the HCL string for data_affinity_regions and the description.
	// dataAffinityRegionsHCL will be something like `["us-central1"]` or `["us-central1", "europe-west1"]`
	updatedContext["data_affinity_regions_hcl"] = dataAffinityRegionsHCL
	updatedContext["description_hcl"] = description

	return acctest.Nprintf(`
# Separate Data Exchange for the Pub/Sub listing to prevent conflicts
resource "google_bigquery_analytics_hub_data_exchange" "listing_pubsub" {
  location         = "US"
  data_exchange_id = "tf_test_pubsub_data_exchange_update_%{random_suffix}"
  display_name     = "tf_test_pubsub_data_exchange_update_%{random_suffix}"
  description      = "Example for pubsub topic source - data exchange%{random_suffix}"
}

# Pub/Sub Topic used as the source for the listing
resource "google_pubsub_topic" "tf_test_pubsub_topic" {
  name = "tf_test_test_pubsub_update_%{random_suffix}"
}

# BigQuery Analytics Hub Listing sourced from the Pub/Sub Topic
resource "google_bigquery_analytics_hub_listing" "listing_pubsub" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing_pubsub.data_exchange_id
  listing_id       = "tf_test_pubsub_listing_update_%{random_suffix}"
  display_name     = "tf_test_pubsub_listing_update_%{random_suffix}"
  description      = "%{description_hcl}" 
  primary_contact  = "test_pubsub_contact@example.com" 

  pubsub_topic {
    topic               = google_pubsub_topic.tf_test_pubsub_topic.id
    data_affinity_regions = %{data_affinity_regions_hcl} 
  }
}
`, updatedContext)
}

func TestAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingMultiregionExample(t *testing.T) {
	t.Parallel()

	bqdataset, err := acctest.AddBigQueryDatasetReplica(envvar.GetTestProjectFromEnv(), "my_listing_example2", "us", "eu")
	if err != nil {
		// If an error occurs, fail the test immediately and log the error.
		t.Fatalf("Failed to create BigQuery dataset and add replica: %v", err)
	}
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"bqdataset":     bqdataset,
	}

	t.Cleanup(func() {
		acctest.CleanupBigQueryDatasetAndReplica(envvar.GetTestProjectFromEnv(), "my_listing_example2", "eu")
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubListingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingMultiregionExample(context),
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

func testAccBigqueryAnalyticsHubListing_bigqueryAnalyticshubListingMultiregionExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing" {
  location         = "us"
  data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
  display_name     = "tf_test_my_data_exchange%{random_suffix}" 
  description      = "example listing for multiregion%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "us"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = "tf_test_my_listing%{random_suffix}"
  display_name     = "tf_test_my_listing%{random_suffix}"
  description      = "example listing for multiregion%{random_suffix}"

  bigquery_dataset {
    dataset = "%{bqdataset}"
    replica_locations = ["eu"]
  }
}
`, context)
}
