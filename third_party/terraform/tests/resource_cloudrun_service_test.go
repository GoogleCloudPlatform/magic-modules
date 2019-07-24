package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCloudRunService_cloudrunServiceUpdate(t *testing.T) {
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
				ResourceName:            "google_cloud_run_service.default",
				ImportStateId:           "us-central1/" + name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "status.0.conditions"},
			},
			{
				Config: testAccCloudRunService_cloudrunServiceUpdate(name, project, "50"),
			},
			{
				ResourceName:            "google_cloud_run_service.default",
				ImportStateId:           "us-central1/" + name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "status.0.conditions"},
			},
		},
	})
}

func testAccCloudRunService_cloudrunServiceUpdate(name, project, concurrency string) string {
	return fmt.Sprintf(`
resource "google_cloud_run_service" "default" {
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
