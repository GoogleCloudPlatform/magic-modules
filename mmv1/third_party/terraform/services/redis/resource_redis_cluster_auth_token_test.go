package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/redis"
)

func TestAccRedisClusterAuthToken_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisClusterAuthToken_basic(context),
			},
			{
				ResourceName:      "google_redis_cluster_auth_token.token-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedisClusterAuthToken_rotation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisClusterAuthToken_rotationStep1(context),
			},
			{
				ResourceName:      "google_redis_cluster_auth_token.token1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            "google_redis_cluster_auth_token.token2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_policy"},
			},
			{
				Config: testAccRedisClusterAuthToken_rotationStep2(context),
			},
		},
	})
}

func testAccRedisClusterAuthToken_basic(context map[string]interface{}) string {
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
`, context)
}

func testAccRedisClusterAuthToken_rotationStep1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster" {
  name                        = "tf-test-redis-cluster-%{random_suffix}"
  shard_count                 = 1
  region                      = "europe-west4"
  authorization_mode          = "AUTH_MODE_TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_redis_cluster_auth_token" "token1" {
  token_auth_user             = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
  token_id                    = "version-1"
}

resource "google_redis_cluster_auth_token" "token2" {
  token_auth_user             = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
  deletion_policy             = "ABANDON"
  depends_on                  = [google_redis_cluster_auth_token.token1]
}
`, context)
}

func testAccRedisClusterAuthToken_rotationStep2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster" {
  name                        = "tf-test-redis-cluster-%{random_suffix}"
  shard_count                 = 1
  region                      = "europe-west4"
  authorization_mode          = "AUTH_MODE_TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_redis_cluster_auth_token" "token2" {
  token_auth_user             = "${google_redis_cluster.cluster.id}/tokenAuthUsers/default"
  deletion_policy             = "ABANDON"
}
`, context)
}
