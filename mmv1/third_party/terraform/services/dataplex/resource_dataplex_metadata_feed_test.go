package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexMetadataFeed_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexMetadataFeedDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexMetadataFeed_full(context),
			},
			{
				ResourceName:            "google_dataplex_metadata_feed.test_metadata_feed",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata_feed_id", "location"},
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
				ImportStateVerifyIgnore: []string{"metadata_feed_id", "location"},
			},
		},
	})
}

func testAccDataplexMetadataFeed_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_dataplex_metadata_feed" "test_metadata_feed" {
  metadata_feed_id = "tf-test-metadata-feed-%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"

  labels = {
    "env" = "test"
  }

  scope {
    projects = ["projects/%{project_name}"]
  }

  filters {
    change_types = ["CREATE", "UPDATE"]
  }

  pubsub_topic = google_pubsub_topic.topic.id
}
`, context)
}

func testAccDataplexMetadataFeed_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "topic" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_pubsub_topic" "topic2" {
  name = "tf-test-topic2-%{random_suffix}"
}

resource "google_dataplex_metadata_feed" "test_metadata_feed" {
  metadata_feed_id = "tf-test-metadata-feed-%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"

  labels = {
    "env" = "prod"
    "new_label" = "value"
  }

  scope {
    projects = ["projects/%{project_name}"]
    organization_level = false
  }

  filters {
    change_types = ["CREATE", "UPDATE", "DELETE"]
  }

  pubsub_topic = google_pubsub_topic.topic2.id
}
`, context)
}
