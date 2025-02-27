package bigqueryanalyticshub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigqueryAnalyticsHubListingSubscription_bigqueryAnalyticshubListingSubscriptionResourceTags_Update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubListingSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubListingSubscription_bigqueryAnalyticshubListingSubscriptionBasicExample(context),
			},
			{
				ResourceName:      "google_bigquery_analytics_hub_listing_subscription.subscription",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigqueryAnalyticsHubListingSubscription_bigqueryAnalyticshubListingSubscriptionResourceTags_update(context),
			},
		},
	})
}

func testAccBigqueryAnalyticsHubListingSubscription_bigqueryAnalyticshubListingSubscriptionResourceTags_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_tags_tag_key" "tag_key1" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key1%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent = google_tags_tag_key.tag_key1.id
  short_name = "tf_test_tag_value1%{random_suffix}"
}

resource "google_bigquery_dataset" "listing" {
  dataset_id                  = "tf_test_my_listing%{random_suffix}"
  friendly_name               = "tf_test_my_listing%{random_suffix}"
  description                 = "example data exchange%{random_suffix}"
  location                    = "US"
}

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

resource "google_bigquery_analytics_hub_listing_subscription" "listing" {
  location         = "US"
  data_exchange_id = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id       = google_bigquery_analytics_hub_listing.listing.listing_id
  destination_dataset {
    description   = "A test subscription"
    friendly_name = "tf_test_my_listing_subscription"
    labels = {
      testing = "123"
    }
    location = "US"
    dataset_reference {
      dataset_id = google_bigquery_dataset.listing.id
      project_id = google_bigquery_dataset.listing.project
    }
    resource_tags = {
    }
  }
}
`, context)
}
