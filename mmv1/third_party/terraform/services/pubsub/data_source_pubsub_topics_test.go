package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/pubsub"
)

func TestAccDataSourceGooglePubsubTopics_basic(t *testing.T) {
	t.Parallel()

	topicName := fmt.Sprintf("tf-test-topics-ds-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubTopicDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubTopicsConfig(topicName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_pubsub_topics.all", "id", fmt.Sprintf("projects/%s/topics", envvar.GetTestProjectFromEnv())),
					// The shared test project holds other topics, so there is no count check and the list must hold this one.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_pubsub_topics.all", "topics.*", map[string]string{
						"name":                       topicName,
						"labels.my-label":            "my-label-value",
						"message_retention_duration": "86600s",
					}),
					resource.TestCheckTypeSetElemAttrPair("data.google_pubsub_topics.all", "topics.*.id", "google_pubsub_topic.foo", "id"),
				),
			},
		},
	})
}

// testAccDataSourceGooglePubsubTopicsConfig reads the project from the topic, so the list is read after the topic exists.
func testAccDataSourceGooglePubsubTopicsConfig(topicName string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%s"
  labels = {
    my-label = "my-label-value"
  }
  message_retention_duration = "86600s"
}

data "google_pubsub_topics" "all" {
  project = google_pubsub_topic.foo.project
}
`, topicName)
}
