package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/redis"
)

func TestAccRedisClusterAuthTokenDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisClusterAuthTokenDatasource(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_redis_cluster_auth_token.default", "google_redis_cluster_auth_token.token-basic"),
				),
			},
		},
	})
}

func testAccRedisClusterAuthTokenDatasource(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster" {
  name                        = "tf-test-redis-cluster-%{random_suffix}"
  shard_count                 = 1
  region                      = "europe-west4"
  deletion_protection_enabled = false
}

resource "google_redis_cluster_token_auth_user" "user-basic" {
  cluster                     = google_redis_cluster.cluster.name
  user_id                     = "tf-test-user-%{random_suffix}"
}

resource "google_redis_cluster_auth_token" "token-basic" {
  token_auth_user             = google_redis_cluster_token_auth_user.user-basic.id
}

data "google_redis_cluster_auth_token" "default" {
  token_auth_user             = google_redis_cluster_token_auth_user.user-basic.id
  token_id                    = google_redis_cluster_auth_token.token-basic.token_id
}
`, context)
}
