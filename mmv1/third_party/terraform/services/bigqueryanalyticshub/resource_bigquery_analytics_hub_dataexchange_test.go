package bigqueryanalyticshub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryAnalyticsHubDataExchange_bigqueryAnalyticshubPublicDataExchangeUpdate(t *testing.T) {
	t.Parallel()

	randString := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"project":          envvar.GetTestProjectFromEnv(),
		"location":         "US",
		"random_suffix":    randString,
		"data_exchange_id": "tf_test_my_data_exchange" + randString,
		"desc":             "description",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubDataExchangeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubDataExchange_bigqueryAnalyticshubPublicDataExchangeExample(context),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_data_exchange.data_exchange",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_exchange_id", "location"},
			},
			{
				Config: testAccBigqueryAnalyticsHubDataExchange_bigqueryAnalyticshubPublicDataExchangeUpdate(context),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_data_exchange.data_exchange",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_exchange_id", "location"},
			},
		},
	})
}

func testAccBigqueryAnalyticsHubDataExchange_bigqueryAnalyticshubPublicDataExchangeUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "data_exchange" {
  location         = "US"
  data_exchange_id = "tf_test_public_data_exchange%{random_suffix}"
  display_name     = "tf_test_public_data_exchange%{random_suffix}"
  description      = "Example for public data exchange%{random_suffix}"
  discovery_type   = "DISCOVERY_TYPE_PRIVATE"
}
`, context)
}
