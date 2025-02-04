package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccStorageBucketIamPolicy_destroy(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketIamPolicy_destroy(),
			},
		},
	})
}

func testAccStorageBucketIamPolicy_destroy() string {
	return fmt.Sprintf(`
resource "google_service_account" "accessor" {
  account_id = "pub-sub-test-service-account"
}

resource "google_storage_bucket" "test_bucket" {
  name          = "sdk-pubsub-test-bucket"
  location      = "US"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  force_destroy = true
}

data "google_iam_policy" "bucket_policy_data" {
  binding {
    role = "roles/storage.admin"

    members = ["serviceAccount:${google_service_account.accessor.email}"]
  }
}

resource "google_storage_bucket_iam_policy" "bucket_policy" {
  bucket      = google_storage_bucket.test_bucket.name
  policy_data = data.google_iam_policy.bucket_policy_data.policy_data
}

resource "google_pubsub_topic" "topic" {
  name = "sdk-pubsub-test-bucket-topic"
}

resource "google_storage_notification" "storage_notification" {
  bucket         = google_storage_bucket.test_bucket.name
  payload_format = "JSON_API_V1"
  topic          = google_pubsub_topic.topic.id

  depends_on = [google_pubsub_topic_iam_policy.topic_policy]
}

data "google_storage_project_service_account" "gcs_account" {}

data "google_iam_policy" "topic_policy_data" {
  binding {
    role = "roles/pubsub.publisher"
    members = [
      "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
    ]
  }
}

resource "google_pubsub_topic_iam_policy" "topic_policy" {
  topic       = google_pubsub_topic.topic.name
  policy_data = data.google_iam_policy.topic_policy_data.policy_data
}
`)
}
