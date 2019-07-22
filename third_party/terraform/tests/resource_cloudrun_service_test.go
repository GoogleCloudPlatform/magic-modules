package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCloudrunService_cloudrunServiceUpdate(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	name := "tftest-cloudrun-" + acctest.RandString(6)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudrunService_cloudrunServiceUpdate(name, project, "10"),
			},
			{
				ResourceName:            "google_cloudrun_service.default",
				ImportStateId:           "us-central1/" + name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
			{
				Config: testAccCloudrunService_cloudrunServiceUpdate(name, project, "50"),
			},
			{
				ResourceName:            "google_cloudrun_service.default",
				ImportStateId:           "us-central1/" + name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
		},
	})
}

func testAccCloudrunService_cloudrunServiceUpdate(name, project, concurrency string) string {
	return fmt.Sprintf(`
resource "google_cloudrun_service" "default" {
  name          = "%s"
  location = "us-central1"

  metadata {
    namespace = "%s"
  }

  spec {
    containers {
	  image = "gcr.io/cloudrun/hello"
	  args = ["arrgs"]
	}
	container_concurrency = %s
  }
}
`, name, project, concurrency)
}
