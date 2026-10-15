package cloudbuild_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCloudBuildGitLabConfig_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"config_id":     "tf-test-gitlab-config" + randomSuffix,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBuildGitLabConfig_full(context),
			},
			{
				ResourceName:            "google_cloudbuild_git_lab_config.gitlab-config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_id", "location"},
			},
			{
				Config: testAccCloudBuildGitLabConfig_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_cloudbuild_git_lab_config.gitlab-config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_cloudbuild_git_lab_config.gitlab-config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_id", "location"},
			},
		},
	})
}

func testAccCloudBuildGitLabConfig_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuild_git_lab_config" "gitlab-config" {
    config_id = "%{config_id}"
    location = "us-central1"
    username = "test-user"
    secrets {
        webhook_secret_version = "projects/myProject/secrets/mysecret/versions/1"
        api_key_version = "projects/myProject/secrets/mysecret/versions/1"
        api_access_token_version = "projects/myProject/secrets/mysecret/versions/1"
        read_access_token_version = "projects/myProject/secrets/mysecret/versions/1"
    }
}
`, context)
}

func testAccCloudBuildGitLabConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloudbuild_git_lab_config" "gitlab-config" {
    config_id = "%{config_id}"
    location = "us-central1"
    username = "updated-user"
    secrets {
        webhook_secret_version = "projects/myProject/secrets/mysecret/versions/1"
        api_key_version = "projects/myProject/secrets/mysecret/versions/1"
        api_access_token_version = "projects/myProject/secrets/mysecret/versions/2"
        read_access_token_version = "projects/myProject/secrets/mysecret/versions/1"
    }
}
`, context)
}
