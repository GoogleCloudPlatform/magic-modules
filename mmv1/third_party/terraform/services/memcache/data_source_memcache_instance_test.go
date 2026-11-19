package memcache_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/memcache"
	"github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
)

func TestAccMemcacheInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"network_name":  servicenetworking.BootstrapSharedServiceNetworkingConnection(t, "vpc-network-1"),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstanceDatasourceConfig(context),
				Check: acctest.CheckDataSourceStateMatchesResourceState(
					"data.google_memcache_instance.default",
					"google_memcache_instance.instance",
				),
			},
		},
	})
}

func testAccMemcacheInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_network" "memcache_network" {
  name = "%{network_name}"
}

resource "google_memcache_instance" "instance" {
  name               = "tf-test-test-instance%{random_suffix}"
  authorized_network = data.google_compute_network.memcache_network.id
  region             = "us-central1"
  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 1
}

data "google_memcache_instance" "default" {
  name   = google_memcache_instance.instance.name
  region = "us-central1"
}
`, context)
}
