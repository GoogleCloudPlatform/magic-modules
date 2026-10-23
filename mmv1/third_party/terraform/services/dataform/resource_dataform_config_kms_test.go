package dataform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dataform"
	_ "github.com/hashicorp/terraform-provider-google/google/services/kms"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

// Covers default_kms_key_name in the GA provider. Unlike TestAccDataformConfig_update, this test
// doesn't use the beta-only google_project_service_identity resource; it grants the KMS role to the
// Dataform service agent by its well-known email instead.
func TestAccDataformConfig_kmsKey(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	withKey := func(key string) map[string]interface{} {
		c := make(map[string]interface{}, len(context)+1)
		for k, v := range context {
			c[k] = v
		}
		if key != "" {
			c["kms_key_attr"] = "default_kms_key_name = " + key
		} else {
			c["kms_key_attr"] = "# KMS key intentionally not set"
		}
		return c
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataformConfig_kmsKey(withKey("google_kms_crypto_key.example_key_1.id")),
			},
			{
				ResourceName:            "google_dataform_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccDataformConfig_kmsKey(withKey("google_kms_crypto_key.example_key_2.id")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataform_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataform_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				Config: testAccDataformConfig_kmsKey(withKey("")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataform_config.config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataform_config.config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccDataformConfig_kmsKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "tf-test-project-2%{random_suffix}"
  project_id      = "tf-test-project-2%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

# Enable Cloud Key KMS API
resource "google_project_service" "cloudkms_api" {
  project            = google_project.project.project_id
  service            = "cloudkms.googleapis.com"
  disable_on_destroy = false
}

# Enable Dataform API
resource "google_project_service" "dataform_api" {
  project            = google_project.project.project_id
  service            = "dataform.googleapis.com"
  disable_on_destroy = false
  depends_on         = [
    google_project_service.cloudkms_api
  ]
}

# Wait for the Dataform service agent to be provisioned and IAM to propagate after API enablement
resource "time_sleep" "wait_for_dataform_api" {
  create_duration = "60s"
  depends_on      = [
    google_project_service.dataform_api
  ]
}

resource "google_kms_key_ring" "keyring" {
  project    = google_project.project.project_id
  name       = "tf-test-example-key-ring%{random_suffix}"
  location   = "us-central1"
  depends_on = [
    time_sleep.wait_for_dataform_api
  ]
}

resource "google_kms_crypto_key" "example_key_1" {
  name     = "tf-test-example-crypto-key-name-1-%{random_suffix}"
  key_ring = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key" "example_key_2" {
  name     = "tf-test-example-crypto-key-name-2-%{random_suffix}"
  key_ring = google_kms_key_ring.keyring.id
}

# Grant the Encrypter/Decrypter role to the Dataform service agent on the KMS keys
resource "google_kms_crypto_key_iam_member" "crypto_key_1_binding" {
  crypto_key_id = google_kms_crypto_key.example_key_1.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${google_project.project.number}@gcp-sa-dataform.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "crypto_key_2_binding" {
  crypto_key_id = google_kms_crypto_key.example_key_2.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${google_project.project.number}@gcp-sa-dataform.iam.gserviceaccount.com"
}

resource "google_dataform_config" "config" {
  region  = "us-central1"
  project = google_project.project.project_id
  %{kms_key_attr}

  # Ensure IAM bindings are created before the Dataform Config
  depends_on = [
    google_kms_crypto_key_iam_member.crypto_key_1_binding,
    google_kms_crypto_key_iam_member.crypto_key_2_binding,
    time_sleep.wait_for_dataform_api
  ]
}
`, context)
}
