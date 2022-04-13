package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudDeployTarget_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataCatalogEntryGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudDeployTarget_deployTargetExample(context),
			},
			{
				ResourceName:            "google_cloud_deploy_target.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region"},
			},
			{
				Config: testAccCloudDeployTarget_deployTargetExample_update(context),
			},
			{
				ResourceName:            "google_cloud_deploy_target.pipeline",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region"},
			},
		},
	})
}

func testAccCloudDeployTarget_deployTargetExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_deploy_target" "pipeline" {
  name          = "tf-test-tf-test%{random_suffix}"
  description   = "Target Cluster"
  annotations = {
    generated-by = "magic-modules"
  }
  labels = {
    env = "dev"
  }
  gke {
    cluster = "${data.google_project.project.id}/locations/us-central1/clusters/prod"
  }
  execution_configs {
    usages = ["RENDER", "DEPLOY"]
    service_account = data.google_app_engine_default_service_account.default.email
  }
}

data "google_project" "project" {
}

data "google_app_engine_default_service_account" "default" {
}
`, context)
}
