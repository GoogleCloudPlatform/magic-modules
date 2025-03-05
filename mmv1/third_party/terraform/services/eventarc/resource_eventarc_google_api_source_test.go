package eventarc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccEventarcGoogleApiSource_update(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"region":          region,
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcGoogleApiSourceDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcGoogleApiSource_full(context),
			},
			{
				ResourceName:            "google_eventarc_google_api_source.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "google_api_source_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccEventarcGoogleApiSource_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_google_api_source.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_google_api_source.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "google_api_source_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccEventarcGoogleApiSource_unset(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_eventarc_google_api_source.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_eventarc_google_api_source.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "google_api_source_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

// Sets up an initial project containing a GoogleApiSource with CMEK connected
// to a MessageBus in the same project
func testAccEventarcGoogleApiSource_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project_1" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_create_project_1" {
  create_duration = "60s"
  depends_on      = [google_project.project_1]
}

resource "google_project_service" "compute_1" {
  project    = google_project.project_1.project_id
  service    = "compute.googleapis.com"
  depends_on = [time_sleep.wait_create_project_1]
}

resource "google_project_service" "servicenetworking_1" {
  project   = google_project.project_1.project_id
  service   = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.compute_1]
}

resource "google_project_service" "kms_1" {
  project    = google_project.project_1.project_id
  service    = "cloudkms.googleapis.com"
  depends_on = [google_project_service.servicenetworking_1]
}

resource "google_project_service" "eventarc_1" {
  project    = google_project.project_1.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.kms_1]
}

resource "time_sleep" "wait_enable_service_1" {
  create_duration = "20s"
  depends_on      = [google_project_service.eventarc_1]
}

resource "google_kms_key_ring" "keyring_1" {
  name       = "keyring"
  location   = "%{region}"
  project    = google_project.project_1.project_id
  depends_on = [google_project_service.kms_1]
}

resource "google_kms_crypto_key" "key_1" {
  name     = "key1"
  key_ring = google_kms_key_ring.keyring_1.id
}

resource "google_project_service_identity" "eventarc_sa_1" {
  service    = "eventarc.googleapis.com"
  project    = google_project.project_1.project_id
  depends_on = [time_sleep.wait_enable_service_1]
}

resource "google_kms_crypto_key_iam_member" "eventarc_sa_keyuser_1" {
  crypto_key_id = google_kms_crypto_key.key_1.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = google_project_service_identity.eventarc_sa_1.member
}

resource "time_sleep" "wait_create_sa_1" {
  create_duration = "20s"
  depends_on      = [google_project_service_identity.eventarc_sa_1, google_kms_crypto_key_iam_member.eventarc_sa_keyuser_1]
}

resource "google_eventarc_message_bus" "message_bus_1" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus%{random_suffix}"
  project        = google_project.project_1.project_id
  depends_on     = [time_sleep.wait_create_sa_1]
}

resource "google_eventarc_google_api_source" "primary" {
  location             = "%{region}"
  google_api_source_id = "tf-test-googleapisource%{random_suffix}"
  project              = google_project.project_1.project_id
  display_name         = "basic google api source"
  destination          = google_eventarc_message_bus.message_bus_1.id
  crypto_key_name      = google_kms_crypto_key.key_1.id
  labels = {
    test_label = "test-eventarc-label"
  }
  annotations = {
    test_annotation = "test-eventarc-annotation"
  }
  logging_config {
    log_severity = "DEBUG"
  }
}
`, context)
}

// Updates all possible fields in the GoogleApiSource, including setting a new
// CMEK key (in the same project) and a new MessageBus (in a different project)
func testAccEventarcGoogleApiSource_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project_1" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute_1" {
  project    = google_project.project_1.project_id
  service    = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking_1" {
  project   = google_project.project_1.project_id
  service   = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.compute_1]
}

resource "google_project_service" "kms_1" {
  project    = google_project.project_1.project_id
  service    = "cloudkms.googleapis.com"
  depends_on = [google_project_service.servicenetworking_1]
}

resource "google_project_service" "eventarc_1" {
  project    = google_project.project_1.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.kms_1]
}

resource "google_kms_key_ring" "keyring_1" {
  name       = "keyring"
  location   = "%{region}"
  project    = google_project.project_1.project_id
  depends_on = [google_project_service.kms_1]
}

resource "google_kms_crypto_key" "key_1" {
  name     = "key1"
  key_ring = google_kms_key_ring.keyring_1.id
}

resource "google_kms_crypto_key" "key_2" {
  name     = "key2"
  key_ring = google_kms_key_ring.keyring_1.id
}

resource "google_project_service_identity" "eventarc_sa_1" {
  service    = "eventarc.googleapis.com"
  project    = google_project.project_1.project_id
}

resource "google_kms_crypto_key_iam_member" "eventarc_sa_keyuser_1" {
  crypto_key_id = google_kms_crypto_key.key_1.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = google_project_service_identity.eventarc_sa_1.member
}

resource "google_kms_crypto_key_iam_member" "eventarc_sa_keyuser_2" {
  crypto_key_id = google_kms_crypto_key.key_2.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = google_project_service_identity.eventarc_sa_1.member
}

resource "time_sleep" "wait_cmek_2" {
  create_duration = "20s"
  depends_on      = [google_kms_crypto_key_iam_member.eventarc_sa_keyuser_2]
}

# Create a separate project to contain another MessageBus.
resource "google_project" "project_2" {
  project_id      = "tf-test2%{random_suffix}"
  name            = "tf-test2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_create_project_2" {
  create_duration = "60s"
  depends_on      = [google_project.project_2]
}

resource "google_project_service" "eventarc_2" {
  project    = google_project.project_2.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [time_sleep.wait_create_project_2]
}

resource "time_sleep" "wait_enable_service_2" {
  create_duration = "20s"
  depends_on      = [google_project_service.eventarc_2]
}

resource "google_project_service_identity" "eventarc_sa_2" {
  project    = google_project.project_2.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [time_sleep.wait_enable_service_2]
}

resource "time_sleep" "wait_create_sa_2" {
  create_duration = "20s"
  depends_on      = [google_project_service_identity.eventarc_sa_2]
}

resource "google_eventarc_message_bus" "message_bus_2" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus2%{random_suffix}"
  project        = google_project.project_2.project_id
  depends_on     = [time_sleep.wait_create_sa_2]
}

resource "google_eventarc_google_api_source" "primary" {
  location             = "%{region}"
  google_api_source_id = "tf-test-googleapisource%{random_suffix}"
  project              = google_project.project_1.project_id
  display_name         = "updated google api source"
  destination          = google_eventarc_message_bus.message_bus_2.id
  crypto_key_name      = google_kms_crypto_key.key_2.id
  labels = {
    updated_label = "updated-test-eventarc-label"
  }
  annotations = {
    updated_test_annotation = "updated-test-eventarc-annotation"
  }
  logging_config {
    log_severity = "ALERT"
  }
  depends_on = [time_sleep.wait_cmek_2]
}
`, context)
}

// Unsets as many fields as possible in the GoogleApiSource.
func testAccEventarcGoogleApiSource_unset(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project_1" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "compute_1" {
  project    = google_project.project_1.project_id
  service    = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking_1" {
  project   = google_project.project_1.project_id
  service   = "servicenetworking.googleapis.com"
  depends_on = [google_project_service.compute_1]
}

resource "google_project_service" "kms_1" {
  project    = google_project.project_1.project_id
  service    = "cloudkms.googleapis.com"
  depends_on = [google_project_service.servicenetworking_1]
}

resource "google_project_service" "eventarc_1" {
  project    = google_project.project_1.project_id
  service    = "eventarc.googleapis.com"
  depends_on = [google_project_service.kms_1]
}

resource "google_kms_key_ring" "keyring_1" {
  name       = "keyring"
  location   = "%{region}"
  project    = google_project.project_1.project_id
  depends_on = [google_project_service.kms_1]
}

resource "google_kms_crypto_key" "key_1" {
  name     = "key1"
  key_ring = google_kms_key_ring.keyring_1.id
}

resource "google_kms_crypto_key" "key_2" {
  name     = "key2"
  key_ring = google_kms_key_ring.keyring_1.id
}

resource "google_project_service_identity" "eventarc_sa_1" {
  service    = "eventarc.googleapis.com"
  project    = google_project.project_1.project_id
}

resource "google_kms_crypto_key_iam_member" "eventarc_sa_keyuser_1" {
  crypto_key_id = google_kms_crypto_key.key_1.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = google_project_service_identity.eventarc_sa_1.member
}

resource "google_kms_crypto_key_iam_member" "eventarc_sa_keyuser_2" {
  crypto_key_id = google_kms_crypto_key.key_2.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = google_project_service_identity.eventarc_sa_1.member
}

resource "google_project" "project_2" {
  project_id      = "tf-test2%{random_suffix}"
  name            = "tf-test2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "eventarc_2" {
  project    = google_project.project_2.project_id
  service    = "eventarc.googleapis.com"
}

resource "google_project_service_identity" "eventarc_sa_2" {
  project    = google_project.project_2.project_id
  service    = "eventarc.googleapis.com"
}

resource "google_eventarc_message_bus" "message_bus_2" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus2%{random_suffix}"
  project        = google_project.project_2.project_id
}

resource "google_eventarc_google_api_source" "primary" {
  location             = "%{region}"
  google_api_source_id = "tf-test-googleapisource%{random_suffix}"
  project              = google_project.project_1.project_id
  destination          = google_eventarc_message_bus.message_bus_2.id
}
`, context)
}

func testAccCheckEventarcGoogleApiSourceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_google_api_source" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{EventarcBasePath}}projects/{{project}}/locations/{{location}}/googleApiSources/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("EventarcGoogleApiSource still exists at %s", url)
			}
		}

		return nil
	}
}
