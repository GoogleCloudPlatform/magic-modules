package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/pubsub"
)

func TestAccDataSourceGooglePubsubSubscriptions_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-subscriptions-ds-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPubsubSubscriptionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGooglePubsubSubscriptionsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_pubsub_subscriptions.all", "id", fmt.Sprintf("projects/%s/subscriptions", envvar.GetTestProjectFromEnv())),
					// The shared test project holds other subscriptions, so there is no count check and the list must hold this one.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_pubsub_subscriptions.all", "subscriptions.*", map[string]string{
						"name":                        name,
						"labels.my-label":             "my-label-value",
						"ack_deadline_seconds":        "20",
						"push_config.0.push_endpoint": "https://example.com/push",
					}),
					resource.TestCheckTypeSetElemAttrPair("data.google_pubsub_subscriptions.all", "subscriptions.*.id", "google_pubsub_subscription.foo", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.google_pubsub_subscriptions.all", "subscriptions.*.topic", "google_pubsub_topic.foo", "id"),
				),
			},
		},
	})
}

// testAccDataSourceGooglePubsubSubscriptionsConfig reads the project from the subscription, so the list is read after the subscription exists.
func testAccDataSourceGooglePubsubSubscriptionsConfig(name string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
  name = "%[1]s"
}

resource "google_pubsub_subscription" "foo" {
  name  = "%[1]s"
  topic = google_pubsub_topic.foo.id
  labels = {
    my-label = "my-label-value"
  }
  ack_deadline_seconds = 20
  push_config {
    push_endpoint = "https://example.com/push"
  }
}

data "google_pubsub_subscriptions" "all" {
  project = google_pubsub_subscription.foo.project
}
`, name)
}
