package google_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	google "internal/terraform-provider-google"
)

func TestAccDataSourceGooglePubsubTopic_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": google.RandString(t, 10),
	}

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubTopic_basic(context),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceState("data.google_pubsub_topic.foo", "google_pubsub_topic.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGooglePubsubTopic_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": google.RandString(t, 10),
	}

	google.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { google.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: google.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubTopic_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					google.CheckDataSourceStateMatchesResourceState("data.google_pubsub_topic.foo", "google_pubsub_topic.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGooglePubsubTopic_basic(context map[string]interface{}) string {
	return google.Nprintf(`
resource "google_pubsub_topic" "foo" {
  name     = "tf-test-pubsub-%{random_suffix}"
}

data "google_pubsub_topic" "foo" {
  name     = google_pubsub_topic.foo.name
  project  = google_pubsub_topic.foo.project
}
`, context)
}

func testAccDataSourceGooglePubsubTopic_optionalProject(context map[string]interface{}) string {
	return google.Nprintf(`
resource "google_pubsub_topic" "foo" {
  name     = "tf-test-pubsub-%{random_suffix}"
}

data "google_pubsub_topic" "foo" {
  name     = google_pubsub_topic.foo.name
}
`, context)
}
