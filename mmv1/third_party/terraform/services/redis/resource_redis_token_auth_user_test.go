package redis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/redis"
)

func TestAccRedisTokenAuthUser_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccRedisTokenAuthUser_basic(context),
			},
			{
				ResourceName:      "google_redis_token_auth_user.user-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRedisTokenAuthUser_basic(context map[string]interface{}) string {
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
`, context)
}
