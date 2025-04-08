package memcache_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemcacheInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemcacheInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemcacheInstanceDatasourceConfig(context),
			},
		},
	})
}

func testAccMemcacheInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memcache_instance" "instance" {
  name = "test-instance"
  node_config {
    cpu_count      = 1
    memory_size_mb = 1024
  }
  node_count = 1
}

data "google_memcache_instance" "default" {
  instance_id                 = google_memcache_instance.instance.name
  location                    = "us-central1"
}
`, context)
}
