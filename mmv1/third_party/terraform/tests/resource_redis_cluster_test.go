package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Validate that replica count is updated for the cluster
func TestAccRedisCluster_updateReplicaCount(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with replica count 1
				Config: createOrUpdateRedisCluster(name /* replicaCount = */, 1 /* shardCount = */, 3, true),
			},
			{
				ResourceName:      "google_redis_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update replica count to 2
				Config: createOrUpdateRedisCluster(name /* replicaCount = */, 2 /* shardCount = */, 3, true),
			},
			{
				ResourceName:      "google_redis_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(name /* replicaCount = */, 2 /* shardCount = */, 3, false),
			},
		},
	})
}

// Validate that shard count is updated for the cluster
func TestAccRedisCluster_updateShardCount(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckRedisClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// create cluster with shard count 3
				Config: createOrUpdateRedisCluster(name /* replicaCount = */, 1 /* shardCount = */, 3, true),
			},
			{
				ResourceName:      "google_redis_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update shard count to 5
				Config: createOrUpdateRedisCluster(name /* replicaCount = */, 1 /* shardCount = */, 5, true),
			},
			{
				ResourceName:      "google_redis_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// clean up the resource
				Config: createOrUpdateRedisCluster(name /* replicaCount = */, 1 /* shardCount = */, 5, false),
			},
		},
	})
}

func createOrUpdateRedisCluster(name string, replicaCount int, shardCount int, preventDestroy bool) string {
	lifecycleBlock := ""
	if preventDestroy {
		lifecycleBlock = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return fmt.Sprintf(`
resource "google_redis_cluster" "test" {
        name           = "%s"
	replica_count = %d
	shard_count = %d
  region         = "us-central1"
	psc_configs {
			network = "projects/${data.google_project.project.number}/global/networks/default"
	}
	%s
}

data "google_project" "project" {
}
`, name, replicaCount, shardCount, lifecycleBlock)
}
