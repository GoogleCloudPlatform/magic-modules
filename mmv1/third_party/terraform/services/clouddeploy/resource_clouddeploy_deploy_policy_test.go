package clouddeploy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccClouddeployDeployPolicy_update(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckClouddeployDeployPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccClouddeployDeployPolicy_basic(),
			},
			{
				ResourceName:            "google_clouddeploy_deploy_policy.deploy-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "annotations", "labels", "terraform_labels"},
			},
			{
				Config: testAccClouddeployDeployPolicy_update(),
			},
			{
				ResourceName:            "google_clouddeploy_custom_target_type.custom-target-type",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "annotations", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccClouddeployDeployPolicy_basic() string {
	return `
resource "google_clouddeploy_deploy_policy" "deploy-policy" {
    location = "us-central1"
    name = "tf-test-my-deploy-policy"
    description = "My deploy policy"
    selectors {
      delivery_pipeline {
        id = "*"
      }
    }
    rules {
      rollout_restriction {
        id = "holidayfreeze"
        time_windows {
          time_zone = "America/New_York"
          one_time_windows {
            start_date = "2024-12-24"
            start_time = "00:00"
            end_date = "2024-12-27"
            end_time = "09:00"
          }
        }
      }
    }
}
`
}

func testAccClouddeployDeployPolicy_update(context map[string]interface{}) string {
	return `
resource "google_clouddeploy_deploy_policy" "deploy-policy" {
    location = "us-central1"
    name = "tf-test-my-deploy-policy"
    description = "My deploy policy"
    selectors {
      delivery_pipeline {
        id = "mypipeline"
      }
    }
    rules {
      rollout_restriction {
        id = "weekendfreeze"
        invokers = ["USER"]
        actions = ["CREATE", "APPROVE"]
        time_windows {
          time_zone = "America/New_York"
          weekly_windows {
            days_of_week = ["SATURDAY", "SUNDAY"]
            start_time = "00:00"
            end_time = "24:00"
          }
        }
      }
    }
    rules {
      rollout_restriction {
        id = "norolloutsMondays"
        time_windows {
          time_zone = "America/New_York"
          weekly_windows {
            days_of_week = ["MONDAY"]
          }
        }
      }
    }
}
`
}
