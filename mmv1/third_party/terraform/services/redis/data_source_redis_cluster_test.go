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
resource "google_redis_cluster" "cluster" {
  name                           = "tf-test-redis-cluster-%{random_suffix}"
  shard_count                    = 1
  region                         = "us-central1"
  deletion_protection_enabled    = false 
  
}   

data "google_redis_cluster" "default" {
  name   = google_redis_cluster.cluster.name
  region = "us-central1"
}
`, context)
}
