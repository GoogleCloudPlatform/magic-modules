package eventarc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventarcTrigger_channel(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	key := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-trigger-key")

	context := map[string]interface{}{
		"region":          region,
		"project_name":    envvar.GetTestProjectFromEnv(),
		"project_number":  envvar.GetTestProjectNumberFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"key_name":        key.CryptoKey.Name,
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_createTriggerWithChannelName(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEventarcTrigger_HttpDest(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()

	testNetworkName := acctest.BootstrapSharedTestNetwork(t, "attachment-network")
	subnetName := acctest.BootstrapSubnet(t, "tf-test-subnet", testNetworkName)
	networkAttachmentName := acctest.BootstrapNetworkAttachment(t, "tf-test-attachment", subnetName)

	// Need to have the full network attachment name in the format project/{project_id}/regions/{region_id}/networkAttachments/{networkAttachmentName}
	fullFormNetworkAttachmentName := fmt.Sprintf("projects/%s/regions/%s/networkAttachments/%s", envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv(), networkAttachmentName)

	context := map[string]interface{}{
		"region":             region,
		"project_name":       envvar.GetTestProjectFromEnv(),
		"project_number":     envvar.GetTestProjectNumberFromEnv(),
		"service_account":    envvar.GetTestServiceAccountFromEnv(t),
		"network_attachment": fullFormNetworkAttachmentName,
		"random_suffix":      acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_createTriggerWithHttpDest(context),
			},
			{
				ResourceName:      "google_eventarc_trigger.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcTrigger_createTriggerWithChannelName(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key_member" {
	crypto_key_id = "%{key_name}"
	role      = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	member = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_channel" "test_channel" {
	location = "%{region}"
	name     = "tf-test-channel%{random_suffix}"
	crypto_key_name = "%{key_name}"
	third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
	depends_on = [google_kms_crypto_key_iam_member.key_member]
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "%{region}"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

resource "google_eventarc_trigger" "primary" {
	name = "tf-test-trigger%{random_suffix}"
	location = "%{region}"
	matching_criteria {
		attribute = "type"
		value = "datadog.v1.alert"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "%{region}"
		}
	}
	service_account = "%{service_account}"

	channel = "projects/%{project_name}/locations/%{region}/channels/${google_eventarc_channel.test_channel.name}"

    depends_on = [google_cloud_run_service.default,google_eventarc_channel.test_channel]
}
`, context)
}

func testAccEventarcTrigger_createTriggerWithHttpDest(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-trigger%{random_suffix}"
	location = "%{region}"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		http_endpoint {
			uri = "http://10.10.10.8:80/route"
		}
                network_config {
                        network_attachment = "%{network_attachment}"
                }

	}
	service_account = "%{service_account}"

}
`, context)
}

func TestAccEventarcTrigger_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccEventarcTrigger_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccEventarcTrigger_BasicHandWrittenUpdate1(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "europe-west1"
		}
	}
	labels = {
		foo = "bar"
	}
}

resource "google_pubsub_topic" "foo" {
	name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

`, context)
}

func testAccEventarcTrigger_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default.name
			region = "europe-west1"
		}
	}
	transport {
		pubsub {
			topic = google_pubsub_topic.foo.id
		}
	}
}

resource "google_pubsub_topic" "foo" {
	name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

resource "google_cloud_run_service" "default2" {
	name     = "tf-test-eventarc-service%{random_suffix}2"
	location = "europe-north1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

`, context)
}

func testAccEventarcTrigger_BasicHandWrittenUpdate1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "europe-west1"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		cloud_run_service {
			service = google_cloud_run_service.default2.name
			region = "europe-north1"
		}
	}
	transport {
		pubsub {
			topic = google_pubsub_topic.foo.id
		}
	}
	labels = {
		foo = "bar"
	}
	service_account = google_service_account.eventarc-sa.email
}

resource "google_service_account" "eventarc-sa" {
	account_id   = "tf-test-sa%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_pubsub_topic" "foo" {
	name = "tf-test-topic%{random_suffix}"
}

resource "google_cloud_run_service" "default" {
	name     = "tf-test-eventarc-service%{random_suffix}"
	location = "europe-west1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

resource "google_cloud_run_service" "default2" {
	name     = "tf-test-eventarc-service%{random_suffix}2"
	location = "europe-north1"

	metadata {
		namespace = "%{project_name}"
	}

	template {
		spec {
			containers {
				image = "gcr.io/cloudrun/hello"
				ports {
					container_port = 8080
				}
			}
			container_concurrency = 50
			timeout_seconds = 100
		}
	}

	traffic {
		percent         = 100
		latest_revision = true
	}
}

`, context)
}

func TestAccEventarcTrigger_workflow(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    envvar.GetTestProjectFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(t),
		"region":          envvar.GetTestRegionFromEnv(),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcTrigger_workflow(context),
			},
			{
				ResourceName:            "google_eventarc_trigger.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccEventarcTrigger_workflow(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_trigger" "primary" {
	name = "tf-test-name%{random_suffix}"
	location = "%{region}"
	matching_criteria {
		attribute = "type"
		value = "google.cloud.pubsub.topic.v1.messagePublished"
	}
	destination {
		workflow = google_workflows_workflow.example.id
	}
	service_account = "%{service_account}"
}

resource "google_workflows_workflow" "example" {
  name = "tf-test-eventarc-workflow%{random_suffix}"
  deletion_protection = false
  region = "%{region}"
  source_contents = <<-EOF
  # This is a sample workflow, feel free to replace it with your source code
  #
  # This workflow does the following:
  # - reads current time and date information from an external API and stores
  #   the response in CurrentDateTime variable
  # - retrieves a list of Wikipedia articles related to the day of the week
  #   from CurrentDateTime
  # - returns the list of articles as an output of the workflow
  # FYI, In terraform you need to escape the $$ or it will cause errors.

  - getCurrentTime:
      call: http.get
      args:
          url: $${sys.get_env("url")}
      result: CurrentDateTime
  - readWikipedia:
      call: http.get
      args:
          url: https://en.wikipedia.org/w/api.php
          query:
              action: opensearch
              search: $${CurrentDateTime.body.dayOfTheWeek}
      result: WikiResult
  - returnOutput:
      return: $${WikiResult.body[1]}
EOF
}
`, context)
}
