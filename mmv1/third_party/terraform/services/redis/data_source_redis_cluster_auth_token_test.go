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
					acctest.CheckDataSourceStateMatchesResourceState("data.google_redis_cluster_auth_token.with_token_id", "google_redis_cluster_auth_token.token-basic"),
					resource.TestCheckResourceAttrSet("data.google_redis_cluster_auth_token.default", "name"),
					resource.TestCheckResourceAttrSet("data.google_redis_cluster_auth_token.default", "token_id"),
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
  authorization_mode          = "AUTH_MODE_TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_redis_cluster_auth_token" "token-basic" {
  token_auth_user             = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
}

data "google_redis_cluster_auth_token" "with_token_id" {
  token_auth_user             = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
  token_id                    = google_redis_cluster_auth_token.token-basic.token_id
}

data "google_redis_cluster_auth_token" "default" {
  token_auth_user             = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
  depends_on                  = [google_redis_cluster_auth_token.token-basic]
}
`, context)
}
