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
				ResourceName:      "google_redis_cluster_auth_token.token-default",
				ImportState:       true,
				ImportStateVerify: true,
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

resource "google_redis_cluster_auth_token" "token-default" {
  cluster                     = google_redis_cluster.cluster.name
  region                      = "europe-west4"
  user_id                     = "default"
}

resource "google_redis_cluster_token_auth_user" "user-basic" {
  cluster                     = google_redis_cluster.cluster.name
  region                      = "europe-west4"
  user_id                     = "tf-test-user-%{random_suffix}"
}

resource "google_redis_cluster_auth_token" "token-basic" {
  cluster                     = google_redis_cluster.cluster.name
  region                      = "europe-west4"
  user_id                     = google_redis_cluster_token_auth_user.user-basic.user_id
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

resource "google_redis_cluster_token_auth_user" "user" {
  cluster                     = google_redis_cluster.cluster.name
  user_id                     = "tf-test-user-%{random_suffix}"
  region                      = "europe-west4"
}

resource "google_redis_cluster_auth_token" "token1" {
  cluster                     = google_redis_cluster.cluster.name
  user_id                     = google_redis_cluster_token_auth_user.user.user_id
  token_id                    = "version-1"
  region                      = "europe-west4"
}

resource "google_redis_cluster_auth_token" "token2" {
  cluster                     = google_redis_cluster.cluster.name
  user_id                     = google_redis_cluster_token_auth_user.user.user_id
  deletion_policy             = "ABANDON"
  region                      = "europe-west4"
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

resource "google_redis_cluster_token_auth_user" "user" {
  cluster                     = google_redis_cluster.cluster.name
  user_id                     = "tf-test-user-%{random_suffix}"
  region                      = "europe-west4"
}

resource "google_redis_cluster_auth_token" "token2" {
  cluster                     = google_redis_cluster.cluster.name
  user_id                     = google_redis_cluster_token_auth_user.user.user_id
  deletion_policy             = "ABANDON"
  region                      = "europe-west4"
}
`, context)
}
