package bigqueryanalyticshub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryAnalyticsHubSubscription_bigqueryAnalyticshubSubscription_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		CheckDestroy: testAccCheckBigqueryAnalyticsHubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubSubscription_bigqueryAnalyticshubSubscription_basic(context),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_subscription.subscription",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_exchange_id", "listing_id", "destination_dataset", "subscription_id", "project"},
			},
			{
				Config: testAccBigqueryAnalyticsHubSubscription_bigqueryAnalyticshubSubscription_update(context),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_subscription.subscription",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "data_exchange_id", "listing_id", "destination_dataset", "subscription_id", "project"},
			},
		},
	})
}

func testAccBigqueryAnalyticsHubSubscription_bigqueryAnalyticshubSubscription_basic(context map[string]interface{}) string {
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
  description      = "example data exchange%{random_suffix}"

  bigquery_dataset {
    dataset = google_bigquery_dataset.listing.id
  }
}

resource "google_bigquery_dataset" "listing" {
  dataset_id       = "tf_test_my_listing%{random_suffix}"
  friendly_name    = "tf_test_my_listing%{random_suffix}"
  description      = ""
  location         = "US"
}

resource "google_bigquery_analytics_hub_subscription" "subscription" {
  location            = "US"
  data_exchange_id    = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id          = google_bigquery_analytics_hub_listing.listing.listing_id
  destination_dataset {
    dataset_reference {
      dataset_id = "tf_test_my_subscription_dataset%{random_suffix}"
      project_id = "%{project}"
    }
    location = "US"
  }
}
`, context)
}

func testAccBigqueryAnalyticsHubSubscription_bigqueryAnalyticshubSubscription_update(context map[string]interface{}) string {
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
	description      = "example data exchange%{random_suffix}"

	bigquery_dataset {
		dataset = google_bigquery_dataset.listing.id
	}
}

resource "google_bigquery_dataset" "listing" {
	dataset_id       = "tf_test_my_listing%{random_suffix}"
	friendly_name    = "tf_test_my_listing%{random_suffix}"
	description      = ""
	location         = "US"
}

resource "google_bigquery_analytics_hub_subscription" "subscription" {
	location            = "US"
  data_exchange_id    = google_bigquery_analytics_hub_data_exchange.listing.data_exchange_id
  listing_id          = google_bigquery_analytics_hub_listing.listing.listing_id

  destination_dataset {
		description = "A new description"
		friendly_name = "A new name"
    dataset_reference {
			dataset_id = "tf_test_my_subscription_dataset%{random_suffix}"
      project_id = "%{project}"
    }
    location = "US"
  }
}
`, context)
}
