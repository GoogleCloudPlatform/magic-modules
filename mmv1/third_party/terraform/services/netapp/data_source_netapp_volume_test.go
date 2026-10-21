package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/netapp"
	"github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
)

func TestAccNetappVolumeDatasource_basic(t *testing.T) {
	t.Parallel()

	location := envvar.GetTestRegionFromEnv()
	if location == "" {
		location = "us-central1"
	}

	context := map[string]interface{}{
		"network_name":  servicenetworking.BootstrapSharedServiceNetworkingConnection(t, "gcnv-network-config-3", servicenetworking.ServiceNetworkWithParentService("netapp.servicenetworking.goog")),
		"location":      location,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappVolumeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappVolumeDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_netapp_volume.default",
						"google_netapp_volume.test_volume",
					),
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_netapp_volume.default_provider_defaults",
						"google_netapp_volume.test_volume",
					),
				),
			},
		},
	})
}

func testAccNetappVolumeDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "default" {
  name = "%{network_name}"
}

resource "google_netapp_storage_pool" "default" {
  name          = "tf-test-pool%{random_suffix}"
  location      = "%{location}"
  service_level = "PREMIUM"
  capacity_gib  = "2048"
  network       = data.google_compute_network.default.id
}

resource "google_netapp_volume" "test_volume" {
  location     = "%{location}"
  name         = "tf-test-volume%{random_suffix}"
  capacity_gib = "100"
  share_name   = "tf-test-volume%{random_suffix}"
  storage_pool = google_netapp_storage_pool.default.name
  protocols    = ["NFSV3"]
  description  = "test volume for datasource"
  labels = {
    foo = "bar"
  }
}

data "google_netapp_volume" "default" {
  name     = google_netapp_volume.test_volume.name
  location = google_netapp_volume.test_volume.location
  project  = google_netapp_volume.test_volume.project
}

data "google_netapp_volume" "default_provider_defaults" {
  name = google_netapp_volume.test_volume.name
}
`, context)
}
