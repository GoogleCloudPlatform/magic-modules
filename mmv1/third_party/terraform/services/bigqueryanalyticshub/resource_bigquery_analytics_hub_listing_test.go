package bigqueryanalyticshub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
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
			},
		},
	})
}

func TestAccBigqueryAnalyticsHubListing_logLinkedDatasetQueryUserEmailCreate(t *testing.T) {
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
				Config: testAccBigqueryAnalyticsHubListing_logLinkedDatasetQueryUserEmailExample(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_analytics_hub_listing.listing", "log_linked_dataset_query_user_email", "true"),
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

func testAccBigqueryAnalyticsHubListing_logLinkedDatasetQueryUserEmailExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "listing_log_email" {
  location         = "US"
  data_exchange_id = "tf_test_log_email_de%{random_suffix}"
  display_name     = "tf_test_log_email_de%{random_suffix}"
  description      = "Data exchange for log email test%{random_suffix}"
}

resource "google_bigquery_analytics_hub_listing" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing_log_email.data_exchange_id
  listing_id       = "tf_test_log_email_listing%{random_suffix}"
  display_name     = "tf_test_log_email_listing%{random_suffix}"
  description      = "Listing with log email enabled%{random_suffix}"
  log_linked_dataset_query_user_email = true

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing_log_email.id
  }
}

resource "google_bigquery_dataset" "listing_log_email" {
  dataset_id                  = "tf_test_log_email_ds%{random_suffix}"
  friendly_name               = "tf_test_log_email_ds%{random_suffix}"
  description                 = "Dataset for log email test%{random_suffix}"
  location                    = "US"
}
`, context)
}
