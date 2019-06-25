package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPubsubTopic_update(t *testing.T) {
	t.Parallel()

	topic := fmt.Sprintf("tf-test-topic-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPubsubTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_update(topic, "foo", "bar"),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPubsubTopic_update(topic, "wibble", "wobble"),
			},
			{
				ResourceName:      "google_pubsub_topic.foo",
				ImportStateId:     topic,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPubsubTopic_cmek(t *testing.T) {
	t.Parallel()

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	topicName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPubsubTopic_cmek(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, topicName),
			},
			{
				ResourceName:      "google_pubsub_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Use a separate TestStep rather than a CheckDestroy because we need the project to still exist.
			{
				Config: testAccPubsubTopic_removed(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
				Check:  testAccCheckPubsubTopicWasRemovedFromState("google_pubsub_topic.topic"),
			},
		},
	})
}

func testAccPubsubTopic_update(topic, key, value string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "foo" {
	name = "%s"
	labels = {
		%s = "%s"
	}
}
`, topic, key, value)
}

// This test runs in its own project, otherwise the test project would start to get filled
// with undeletable resources
func testAccPubsubTopic_cmek(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, topicName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudkms.googleapis.com",
    "pubsub.googleapis.com",
  ]
}

resource "google_kms_key_ring" "key_ring" {
  project  = "${google_project_services.acceptance.project}"
  name     = "%s"
  location = "global"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = "${google_kms_key_ring.key_ring.self_link}"
}

resource "google_project_iam_member" "svc-acct" {
  project = "${google_project_services.acceptance.project}"
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${google_project.acceptance.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

resource "google_pubsub_topic" "topic" {
  name         = "%s"
  project      = "${google_project_iam_member.svc-acct.project}"
  kms_key_name = "${google_kms_crypto_key.crypto_key.self_link}"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, topicName)
}

func testAccPubsubTopic_removed(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_services" "acceptance" {
  project = "${google_project.acceptance.project_id}"

  services = [
    "cloudkms.googleapis.com",
    "pubsub.googleapis.com",
  ]
}

resource "google_kms_key_ring" "key_ring" {
  project  = "${google_project_services.acceptance.project}"
  name     = "%s"
  location = "global"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "%s"
  key_ring = "${google_kms_key_ring.key_ring.self_link}"
}

resource "google_project_iam_member" "svc-acct" {
  project = "${google_project_services.acceptance.project}"
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${google_project.acceptance.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}
`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName)
}

func testAccCheckPubsubTopicWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}
