package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccRedisClusterDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisClusterDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_redis_cluster.default", "google_redis_cluster.cluster"),
				),
			},
		},
	})
}

func testAccRedisClusterDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name = "tf-test-network%{random_suffix}"
}

resource "google_compute_global_address" "service_range" {
  name          = "tf-test-address%{random_suffix}"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.network.id
}

resource "google_service_networking_connection" "private_service_connection" {
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.service_range.name]
}

resource "google_redis_cluster" "cluster" {
  name               = "tf-test-cluster%{random_suffix}"
  region             = "us-central1"
  shard_count        = 2
  replica_count      = 1
  transit_encryption_mode = "TRANSIT_ENCRYPTION_MODE_DISABLED"
  authorization_mode = "AUTH_MODE_DISABLED"
  psc_configs {
    network = google_compute_network.network.id
  }
  depends_on = [google_service_networking_connection.private_service_connection]
  deletion_protection_enabled = false
}

data "google_redis_cluster" "default" {
  name   = google_redis_cluster.cluster.name
  region = "us-central1"
}
`, context)
}
