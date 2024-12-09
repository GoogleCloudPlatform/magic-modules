package gkehub2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleGkeHubFeature_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		Providers:    acctest.TestAccProviders,
		CheckDestroy: testAccCheckGoogleGkeHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleGkeHubFeature_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_gke_hub_feature.example", "google_gke_hub_feature.example"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleGkeHubFeature_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_hub_feature" "example" {
  location = "us-central1"
  project  = "%{project}"
  name     = "configmanagement"
}

data "google_gke_hub_feature" "example" {
  location = google_gke_hub_feature.feature.location
  project  = google_gke_hub_feature.feature.project
  name     = google_gke_hub_feature.feature.name
}
`, context)
}
