package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingLogView_loggingLogViewIamPolicyBasicExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingLogView_loggingLogViewIamPolicyBasicExampleUpdate(context),
			},
			{
				ResourceName:            "google_logging_log_view_iam_policy.log_view_iam_policy_0",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "location", "bucket"},
			},
			{
				Config: testAccLoggingLogView_loggingLogViewIamPolicyMultiExamplesUpdate(context),
			},
			{
				ResourceName:            "google_logging_log_view_iam_policy.log_view_iam_policy_1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "location", "bucket"},
			},
		},
	})
}

func testAccLoggingLogView_loggingLogViewIamPolicyBasicExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_project_bucket_config" "logging_log_bucket" {
    project        = "%{project}"
    location       = "global"
    retention_days = 30
    bucket_id      = "_Default"
}

resource "google_logging_log_view" "logging_log_view" {
  name        = "tf-test-view%{random_suffix}"
  bucket      = google_logging_project_bucket_config.logging_log_bucket.id
  description = "An updated logging view configured with Terraform"
  filter      = "SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\""
}

data "google_iam_policy" "iam_policy_0" {
  binding {
    role    = "roles/logging.viewAccessor"
    members = ["user:user@domain.com"]
  }
}

resource "google_logging_log_view_iam_policy" "log_view_iam_policy_0" {
  parent      = google_logging_log_view.logging_log_view.parent
  location    = google_logging_log_view.logging_log_view.location
  bucket      = google_logging_log_view.logging_log_view.bucket
  name        = google_logging_log_view.logging_log_view.name
  policy_data = data.google_iam_policy.iam_policy_0.policy_data
}
`, context)
}

func testAccLoggingLogView_loggingLogViewIamPolicyMultiExamplesUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_project_bucket_config" "logging_log_bucket" {
    project        = "%{project}"
    location       = "global"
    retention_days = 30
    bucket_id      = "_Default"
}

data "google_iam_policy" "iam_policy_0" {
  binding {
    role    = "roles/logging.viewAccessor"
    members = ["user:user@domain.com"]
  }
}

resource "google_logging_log_view" "logging_log_view0" {
  name        = "tf-test-view-0-%{random_suffix}"
  bucket      = google_logging_project_bucket_config.logging_log_bucket.id
  description = "An updated logging view configured with Terraform"
  filter      = "SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\""
}

resource "google_logging_log_view_iam_policy" "log_view_iam_policy_0" {
  parent      = google_logging_log_view.logging_log_view0.parent
  location    = google_logging_log_view.logging_log_view0.location
  bucket      = google_logging_project_bucket_config.logging_log_bucket.bucket_id
  name        = google_logging_log_view.logging_log_view0.name
  policy_data = data.google_iam_policy.iam_policy_0.policy_data
}

resource "google_logging_log_view" "logging_log_view1" {
  name        = "tf-test-view-1-%{random_suffix}"
  bucket      = google_logging_project_bucket_config.logging_log_bucket.id
  description = "An updated logging view configured with Terraform"
  filter      = "SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\""
}

resource "google_logging_log_view_iam_policy" "log_view_iam_policy_1" {
  parent      = google_logging_log_view.logging_log_view1.parent
  location    = google_logging_log_view.logging_log_view1.location
  bucket      = google_logging_project_bucket_config.logging_log_bucket.bucket_id
  name        = google_logging_log_view.logging_log_view1.name
  policy_data = data.google_iam_policy.iam_policy_0.policy_data
}
`, context)
}
