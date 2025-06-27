package parallelstore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccParallelstoreInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParallelstoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParallelstoreInstanceDatasourceConfig(context),
			},
		},
	})
}

func testAccParallelstoreInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parallelstore_instance" "instance" {
  instance_id = "instance%{random_suffix}"
  location = "us-central1-a"
  description = "test instance"
  capacity_gib = 12000
  network = google_compute_network.network.name
  reserved_ip_range = google_compute_global_address.private_ip_alloc.name
  deployment_type = "SCRATCH"
  file_stripe_level = "FILE_STRIPE_LEVEL_MIN"
  directory_stripe_level = "DIRECTORY_STRIPE_LEVEL_MIN"
  labels = {
    test = "value"
  }
  depends_on = [google_service_networking_connection.default]
}

resource "google_compute_network" "network" {
  name                    = "network%{random_suffix}"
  auto_create_subnetworks = true
  mtu = 8896
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name          = "address%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = google_compute_network.network.id
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}

data "google_parallelstore_instance" "default" {
  name = google_parallelstore_instance.instance.instance_id
  location = "us-central1-a"
}
`, context)
}
