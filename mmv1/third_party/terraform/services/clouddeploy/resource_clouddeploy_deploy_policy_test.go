package clouddeploy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccClouddeployDeployPolicy_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployDeployPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployDeployPolicy_basic(context),
			},
			{
				ResourceName:            "google_clouddeploy_deploy_policy.deploy_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccClouddeployDeployPolicy_update(context),
			},
			{
				ResourceName:            "google_clouddeploy_deploy_policy.deploy_policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccClouddeployDeployPolicy_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_deploy_policy" "deploy_policy" {
  name     = "tf-test-cd-policy%{random_suffix}"
  location = "us-central1"
  selectors {
    delivery_pipeline {
      id = "tf-test-cd-pipeline%{random_suffix}"

    }
  }
  rules {
    rollout_restriction {
      id = "rule"
      time_windows {
        time_zone = "America/Los_Angeles"
        weekly_windows {
            start_time ="12:00"
            end_time = "13:00"
        }
      }
    }
  }
}
`, context)
}

func testAccClouddeployDeployPolicy_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_clouddeploy_deploy_policy" "deploy_policy" {
  name     = "tf-test-cd-policy%{random_suffix}"
  location = "us-central1"
  selectors {
    delivery_pipeline {
      id = "tf-test-cd-pipeline%{random_suffix}"
    }
  }
  rules {
    rollout_restriction {
      id = "rule"
      time_windows {
        time_zone = "America/Los_Angeles"
        weekly_windows {
            start_time ="12:00"
            end_time = "13:00"
        }
      }
    }
  }
  rules {
    rollout_restriction {
        id = "rule2"
        invokers = ["USER"] 
        actions = ["CREATE"]
        time_windows {
        time_zone = "America/Los_Angeles"
        weekly_windows {
            start_time ="12:00"
            end_time = "13:00"
            days_of_week = ["MONDAY"]
          }

        one_time_windows {
          start = "2024-12-22 17:00"
          end = "2025-01-02 09:00"
      }
     }
    }
  }
}
`, context)
}
