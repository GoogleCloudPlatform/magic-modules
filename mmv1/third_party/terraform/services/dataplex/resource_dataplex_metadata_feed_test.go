package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataplexMetadataFeed_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexMetadataFeed_basic(context),
			},
			{
				ResourceName:            "google_dataplex_metadata_feed.test_metadata_feed",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "metadata_feed_id", "terraform_labels"},
			},
			{
				Config: testAccDataplexMetadataFeed_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_metadata_feed.test_metadata_feed", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataplex_metadata_feed.test_metadata_feed",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "metadata_feed_id", "terraform_labels"},
			},
		},
	})
}

func TestAccDataplexMetadataFeed_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexMetadataFeed_full(context),
			},
			{
				ResourceName:            "google_dataplex_metadata_feed.test_metadata_feed",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "metadata_feed_id", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexMetadataFeed_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "test_topic" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_project_service_identity" "dataplex_identity" {
  project  = google_pubsub_topic.test_topic.project
  service  = "dataplex.googleapis.com"
}

resource "google_pubsub_topic_iam_member" "dataplex_publisher" {
  topic  = google_pubsub_topic.test_topic.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:${google_project_service_identity.dataplex_identity.email}"
}

resource "google_pubsub_topic_iam_member" "dataplex_viewer" {
  topic  = google_pubsub_topic.test_topic.name
  role   = "roles/pubsub.viewer"
  member = "serviceAccount:${google_project_service_identity.dataplex_identity.email}"
}

resource "google_dataplex_metadata_feed" "test_metadata_feed" {
  metadata_feed_id = "tf-test-metadata-feed%{random_suffix}"
  location = "us-central1"
  scope {
    organization_level = true
  }
  pubsub_topic = google_pubsub_topic.test_topic.id

  depends_on = [
    google_pubsub_topic_iam_member.dataplex_publisher,
    google_pubsub_topic_iam_member.dataplex_viewer,
  ]
}
`, context)
}

func testAccDataplexMetadataFeed_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "test_topic" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_project_service_identity" "dataplex_identity" {
  project  = google_pubsub_topic.test_topic.project
  service  = "dataplex.googleapis.com"
}

resource "google_pubsub_topic_iam_member" "dataplex_publisher" {
  topic  = google_pubsub_topic.test_topic.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:${google_project_service_identity.dataplex_identity.email}"
}

resource "google_pubsub_topic_iam_member" "dataplex_viewer" {
  topic  = google_pubsub_topic.test_topic.name
  role   = "roles/pubsub.viewer"
  member = "serviceAccount:${google_project_service_identity.dataplex_identity.email}"
}

resource "google_dataplex_metadata_feed" "test_metadata_feed" {
  metadata_feed_id = "tf-test-metadata-feed%{random_suffix}"
  location = "us-central1"
  scope {
    organization_level = true
  }
  pubsub_topic = google_pubsub_topic.test_topic.id

  labels = {
    foo = "bar"
  }

  depends_on = [
    google_pubsub_topic_iam_member.dataplex_publisher,
    google_pubsub_topic_iam_member.dataplex_viewer,
  ]
}
`, context)
}

func testAccDataplexMetadataFeed_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "test_topic" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_project_service_identity" "dataplex_identity" {
  project  = google_pubsub_topic.test_topic.project
  service  = "dataplex.googleapis.com"
}

resource "google_pubsub_topic_iam_member" "dataplex_publisher" {
  topic  = google_pubsub_topic.test_topic.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:${google_project_service_identity.dataplex_identity.email}"
}

resource "google_pubsub_topic_iam_member" "dataplex_viewer" {
  topic  = google_pubsub_topic.test_topic.name
  role   = "roles/pubsub.viewer"
  member = "serviceAccount:${google_project_service_identity.dataplex_identity.email}"
}

resource "google_dataplex_entry_group" "test_entry_group" {
  entry_group_id = "tf-test-entry-group%{random_suffix}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "test_entry_type" {
  entry_type_id = "tf-test-entry-type%{random_suffix}"
  location = "us-central1"
}

resource "google_dataplex_aspect_type" "test_aspect_type" {
  aspect_type_id = "tf-test-aspect-type%{random_suffix}"
  location = "us-central1"

  data_classification = "DATA_CLASSIFICATION_UNSPECIFIED"
  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_metadata_feed" "test_metadata_feed" {
  metadata_feed_id = "tf-test-metadata-feed%{random_suffix}"
  location = "us-central1"

  scope {
    projects = ["projects/${google_dataplex_entry_group.test_entry_group.project}"]
    entry_groups = ["projects/${google_dataplex_entry_group.test_entry_group.project}/locations/us-central1/entryGroups/${google_dataplex_entry_group.test_entry_group.entry_group_id}"]
  }

  filters {
    entry_types = ["projects/${google_dataplex_entry_group.test_entry_group.project}/locations/us-central1/entryTypes/${google_dataplex_entry_type.test_entry_type.entry_type_id}"]
    aspect_types = ["projects/${google_dataplex_entry_group.test_entry_group.project}/locations/us-central1/aspectTypes/${google_dataplex_aspect_type.test_aspect_type.aspect_type_id}"]
    change_types = ["CREATE", "UPDATE", "DELETE"]
  }

  pubsub_topic = google_pubsub_topic.test_topic.id

  labels = {
    foo = "bar"
  }

  depends_on = [
    google_pubsub_topic_iam_member.dataplex_publisher,
    google_pubsub_topic_iam_member.dataplex_viewer,
  ]
}
`, context)
}
