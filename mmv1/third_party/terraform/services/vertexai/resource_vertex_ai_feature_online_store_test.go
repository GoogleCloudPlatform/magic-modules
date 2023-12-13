package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVertexAIFeatureOnlineStore_vertexAiFeatureonlinestoreWithBigtable_updated(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIFeatureOnlineStore_vertexAiFeatureonlinestoreWithBigtable_basic(context),
			},
			{
				ResourceName:            "google_vertex_ai_feature_online_store.featureonlinestore_bigtable",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "region", "force_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccVertexAIFeatureOnlineStore_vertexAiFeatureonlinestoreWithBigtableExample_update(context),
			},
			{

				ResourceName:            "google_vertex_ai_feature_online_store.featureonlinestore_bigtable",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "etag", "region", "force_destroy", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccVertexAIFeatureOnlineStore_vertexAiFeatureonlinestoreWithBigtable_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
##FeatureOnlineStore With BigTable
resource "google_vertex_ai_feature_online_store" "featureonlinestore_bigtable" {
  name     = "terraform2%{random_suffix}"
  labels = {
    foo = "bar"
  }
  region   = "us-central1"
  bigtable {
    auto_scaling {
        min_node_count = 1
        max_node_count = 3
        cpu_utilization_target = 50
    }
  }
  force_destroy = true
}
`, context)
}

func testAccVertexAIFeatureOnlineStore_vertexAiFeatureonlinestoreWithBigtableExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
##FeatureOnlineStore With BigTable
resource "google_vertex_ai_feature_online_store" "featureonlinestore_bigtable" {
  name     = "terraform2%{random_suffix}"
  labels = {
    foo1 = "bar1"
  }
  region   = "us-central1"
  bigtable {
    auto_scaling {
        min_node_count = 2
        max_node_count = 4
        cpu_utilization_target = 60
    }
  }
  force_destroy = true
}
`, context)
}
