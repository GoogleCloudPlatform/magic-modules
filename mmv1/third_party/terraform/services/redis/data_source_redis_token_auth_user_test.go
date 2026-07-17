package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/redis"
)

func TestAccRedisTokenAuthUserDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisTokenAuthUserDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_redis_token_auth_user.default", "google_redis_token_auth_user.user-basic"),
				),
			},
		},
	})
}

func testAccRedisTokenAuthUserDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_redis_cluster" "cluster" {
  name                        = "tf-test-redis-cluster-%{random_suffix}"
  shard_count                 = 1
  region                      = "us-central1"
  deletion_protection_enabled = false
}

resource "google_redis_token_auth_user" "user-basic" {
  cluster = google_redis_cluster.cluster.name
  user_id = "tf-test-user-%{random_suffix}"
}

data "google_redis_token_auth_user" "default" {
  cluster = google_redis_cluster.cluster.name
  user_id = google_redis_token_auth_user.user-basic.user_id
}
`, context)
}
