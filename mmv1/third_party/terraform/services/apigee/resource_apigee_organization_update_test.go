package apigee_test

import (
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccApigeeOrganization_update(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	default_context := map[string]interface{}{
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}
	update_context := maps.Clone(default_context)
	update_context["org_description"] = "Updated Apigee Org description."

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeOrganization_full(default_context),
			},
			{
				ResourceName:            "google_apigee_organization.org",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "properties", "retention"},
			},
			{
				Config: testAccApigeeOrganization_update(update_context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_apigee_organization.org", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_apigee_organization.org",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "properties", "retention"},
			},
		},
	})
}

func testAccApigeeOrganization_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  provider = google

  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  provider = google

  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "compute" {
  provider = google

  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  provider = google

  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
}

resource "google_project_service" "kms" {
  provider = google

  project = google_project.project.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_compute_network" "apigee_network" {
  provider = google

  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [google_project_service.compute]
}

resource "google_compute_global_address" "apigee_range" {
  provider = google

  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.apigee_network.id
  project       = google_project.project.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  provider = google

  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_kms_key_ring" "apigee_keyring" {
  provider = google

  name       = "apigee-keyring"
  location   = "us-central1"
  project    = google_project.project.project_id
  depends_on = [google_project_service.kms]
}

resource "google_kms_crypto_key" "apigee_key" {
  provider = google

  name            = "apigee-key"
  key_ring        = google_kms_key_ring.apigee_keyring.id
}

resource "google_project_service_identity" "apigee_sa" {
  provider = google

  project = google_project.project.project_id
  service = google_project_service.apigee.service
}

resource "google_kms_crypto_key_iam_member" "apigee_sa_keyuser" {
  provider = google

  crypto_key_id = google_kms_crypto_key.apigee_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = google_project_service_identity.apigee_sa.member
}

resource "google_apigee_organization" "org" {
  provider = google

  display_name                         = "apigee-org"
  description                          = "Terraform-managed Apigee Org"
  analytics_region                     = "us-central1"
  project_id                           = google_project.project.project_id
  authorized_network                   = google_compute_network.apigee_network.id
  billing_type                         = "EVALUATION"
  runtime_database_encryption_key_name = google_kms_crypto_key.apigee_key.id
  properties {
    property {
      name = "features.hybrid.enabled"
      value = "true"
    }
    property {
      name = "features.mart.connect.enabled"
      value = "true"
    }
  }

  depends_on = [
    google_service_networking_connection.apigee_vpc_connection,
    google_kms_crypto_key_iam_member.apigee_sa_keyuser,
  ]
}
`, context)
}

func testAccApigeeOrganization_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  provider = google

  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "google_project_service" "apigee" {
  provider = google

  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "compute" {
  provider = google

  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  provider = google

  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
}

resource "google_project_service" "kms" {
  provider = google

  project = google_project.project.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_compute_network" "apigee_network" {
  provider = google

  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [google_project_service.compute]
}

resource "google_compute_global_address" "apigee_range" {
  provider = google

  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.apigee_network.id
  project       = google_project.project.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  provider = google

  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_kms_key_ring" "apigee_keyring" {
  provider = google

  name       = "apigee-keyring"
  location   = "us-central1"
  project    = google_project.project.project_id
  depends_on = [google_project_service.kms]
}

resource "google_kms_crypto_key" "apigee_key" {
  provider = google

  name            = "apigee-key"
  key_ring        = google_kms_key_ring.apigee_keyring.id
}

resource "google_project_service_identity" "apigee_sa" {
  provider = google

  project = google_project.project.project_id
  service = google_project_service.apigee.service
}

resource "google_kms_crypto_key_iam_member" "apigee_sa_keyuser" {
  provider = google

  crypto_key_id = google_kms_crypto_key.apigee_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  member = google_project_service_identity.apigee_sa.member
}

resource "google_apigee_organization" "org" {
  provider = google

  display_name                         = "apigee-org"
  description                          = "%{org_description}"
  analytics_region                     = "us-central1"
  project_id                           = google_project.project.project_id
  authorized_network                   = google_compute_network.apigee_network.id
  billing_type                         = "EVALUATION"
  runtime_database_encryption_key_name = google_kms_crypto_key.apigee_key.id
  properties {
    property {
      name = "features.hybrid.enabled"
      value = "true"
    }
    property {
      name = "features.mart.connect.enabled"
      value = "true"
    }
  }

  depends_on = [
    google_service_networking_connection.apigee_vpc_connection,
    google_kms_crypto_key_iam_member.apigee_sa_keyuser,
  ]
}
`, context)
}
