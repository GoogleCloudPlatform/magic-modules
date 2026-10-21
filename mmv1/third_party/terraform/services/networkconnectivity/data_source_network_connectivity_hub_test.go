package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/networkconnectivity"
)

func TestAccDataSourceNetworkConnectivityHub_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityHubDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkConnectivityHub_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_network_connectivity_hub.default", "google_network_connectivity_hub.primary"),
					acctest.CheckDataSourceStateMatchesResourceState("data.google_network_connectivity_hub.default_no_project", "google_network_connectivity_hub.primary"),
				),
			},
		},
	})
}

func testAccDataSourceNetworkConnectivityHub_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub-%{random_suffix}"
  description = "A sample hub"
  labels = {
    label-one = "value-one"
  }
}

data "google_network_connectivity_hub" "default" {
  name    = google_network_connectivity_hub.primary.name
  project = google_network_connectivity_hub.primary.project
}

data "google_network_connectivity_hub" "default_no_project" {
  name = google_network_connectivity_hub.primary.name
}
`, context)
}
