package google_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRedisInstanceDatasource_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:  func() { acctest.TestAccPreCheck(t) },
		Providers: acctest.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceDatasourceConfig(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					CheckDataSourceStateMatchesResourceState("data.google_redis_instance.redis", "google_redis_instance.redis"),
				),
			},
		},
	})
}

func testAccRedisInstanceDatasourceConfig(suffix string) string {
	return fmt.Sprintf(`
resource "google_redis_instance" "redis" {
  name               = "redis-test-%s"
  memory_size_gb     = 1
}

data "google_redis_instance" "redis" {
  name = google_redis_instance.redis.name
}
`, suffix)
}
