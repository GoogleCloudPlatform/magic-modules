package pubsub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccPubsubSubscriptionIamMemberListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	suffix := acctest.RandString(t, 10)
	topic := "tf-test-topic-" + suffix
	subscription := "tf-test-sub-" + suffix
	subscriptionId := fmt.Sprintf("projects/%s/subscriptions/%s", project, subscription)
	account := "tf-test-pubsub-iam-" + suffix
	role := "roles/pubsub.viewer"
	member := "serviceAccount:" + envvar.ServiceAccountCanonicalEmail(account)

	fmt.Printf("\n[Expected------]\n")
	fmt.Printf("project=%q\n", project)
	fmt.Printf("subscription=%q\n", subscription)
	fmt.Printf("role=%q\n", role)
	fmt.Printf("member=%q\n", member)
	fmt.Printf("condition_title=<null>\n\n")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubSubscriptionIamMember(project, topic, account, subscription, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_pubsub_subscription_iam_member.test", "subscription", subscription),
					resource.TestCheckResourceAttr("google_pubsub_subscription_iam_member.test", "role", role),
					resource.TestCheckResourceAttr("google_pubsub_subscription_iam_member.test", "member", member),
				),
			},
			{
				Query:  true,
				Config: testAccPubsubSubscriptionIamMemberListQueryWithFilters(subscription, project, role, member),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("google_pubsub_subscription_iam_member.test", 1),
					querycheck.ExpectIdentity("google_pubsub_subscription_iam_member.test", map[string]knownvalue.Check{
						"subscription":    knownvalue.StringExact(subscriptionId),
						"role":            knownvalue.StringExact(role),
						"member":          knownvalue.StringExact(member),
						"project":         knownvalue.StringExact(project),
						"condition_title": knownvalue.Null(),
					}),
				},
			},
		},
	})
}

func testAccPubsubSubscriptionIamMember(project, topic, account, subscription, role string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "topic" {
  project = "%s"
  name  = "%s"
}

resource "google_service_account" "test-account" {
  project = "%s"
  account_id   = "%s"
  display_name = "Pubsub subscription IAM Testing Account"
}

resource "google_pubsub_subscription" "test-sub" {
    project = "%s"
    name = "%s"
    topic = google_pubsub_topic.topic.id
}
resource "google_pubsub_subscription_iam_member" "test" {
  project= "%s"
  subscription = google_pubsub_subscription.test-sub.name
  role   = "%s"
  member = "serviceAccount:${google_service_account.test-account.email}"
}
`, project, topic, project, account, project, subscription, project, role)
}

func testAccPubsubSubscriptionIamMemberListQueryWithFilters(subscription, project, role, member string) string {
	return fmt.Sprintf(`
list "google_pubsub_subscription_iam_member" "test" {
  provider = google
  include_resource = true

  config {
    subscription = %q
    project = %q
    role   = %q
    member = %q
  }
}
`, subscription, project, role, member)
}
