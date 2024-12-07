package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleGKEHub2Feature_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHub2FeatureDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleGKEHub2Feature_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_gke_hub_feature.feature", "google_gke_hub_feature.feature"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleGKEHub2Feature_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_gke_hub_feature" "feature" {
  name = "servicemesh"
  location = "global"
  project = data.google_project.project.project_id

  depends_on = [time_sleep.wait_for_gkehub_enablement]
}

data "google_gke_hub_feature" "feature" {
  location = google_gke_hub_feature.feature.location
  project = data.google_project.project.project_id
  feature = google_gke_hub_feature.feature.name
}
`, context)
}
