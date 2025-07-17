package bigqueryanalyticshub_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBigqueryAnalyticsHubQueryTemplate_bigqueryAnalyticshubQuerytemplateUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryAnalyticsHubQueryTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryAnalyticsHubQueryTemplate_bigqueryAnalyticshubQuerytemplateBasicExample(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigquery_analytics_hub_query_template.querytemplate", "submit", "false"),
				),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_query_template.querytemplate",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_exchange_id", "location", "query_template_id", "submit"},
			},
			{
				Config: testAccBigqueryAnalyticsHubQueryTemplate_bigqueryAnalyticshubQuerytemplateUpdate(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigquery_analytics_hub_query_template.querytemplate", "submit", "true"),
				),
			},
			{
				ResourceName:            "google_bigquery_analytics_hub_query_template.querytemplate",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_exchange_id", "location", "query_template_id", "submit"},
			},
		},
	})
}

func testAccBigqueryAnalyticsHubQueryTemplate_bigqueryAnalyticshubQuerytemplateUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_analytics_hub_data_exchange" "querytemplate" {
	display_name = "My Audience Data Exchange"
	data_exchange_id = "tf_test_my_data_exchange%{random_suffix}"
	description = "example of query template%{random_suffix}"
	location = "us"
	sharing_environment_config {
	dcr_exchange_config {}
	}
}

resource "google_bigquery_analytics_hub_query_template" "querytemplate" {
	location = "us"
	data_exchange_id = google_bigquery_analytics_hub_data_exchange.querytemplate.data_exchange_id
	query_template_id = "qt1%{random_suffix}"
	display_name = "qt1%{random_suffix}"
	description = "updated example of query template%{random_suffix}"
	primary_contact = "adminupdated@example.com"
	documentation = "This TVF takes a table t1 as input and returns all columns. Useful for basic data pass-through."
	routine {
		routine_type="TABLE_VALUED_FUNCTION"
		definition_body="qt1%{random_suffix}() as (select * from t1updated)"
	}
	submit=true
} 
`, context)
}
