package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCloudIdsEndpoint_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIdsEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testCloudIds_basic(context),
			},
			{
				ResourceName:      "google_cloud_ids_endpoint.terraform-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCloudIds_basic(context map[string]interface{}) string {
    return Nprintf(`
resource "google_cloud_ids_endpoint" "endpoint" {
	name     = "cloud-ids-test-%{random_suffix}"
	location = "us-central1-f"
	network  = "src-net"
    severity = "INFORMATIONAL"
}
`, context)
}
