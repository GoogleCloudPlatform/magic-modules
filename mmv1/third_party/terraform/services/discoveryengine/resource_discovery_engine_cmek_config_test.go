package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_cmek_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cmek_config_id", "location", "project", "set_default"},
			},
			{
				Config: testAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_update(context),
			},
			{
				ResourceName:            "google_discovery_engine_cmek_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cmek_config_id", "location", "project", "set_default"},
			},
		},
	})
}

func testAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_discovery_engine_cmek_config" "default" {
  location            = "us"
  cmek_config_id      = "tf-test-cmek-config-id%{random_suffix}"
  kms_key             = "projects/${data.google_project.project.project_id}/locations/us/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-1"
  kms_key_version     = "projects/${data.google_project.project.project_id}/locations/us/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-1/cryptoKeyVersions/1"
  set_default         = true

  depends_on = [
    google_kms_crypto_key_iam_binding.discoveryengine_cmek_keyuser
  ]
}

resource "google_kms_key_ring" "key_ring" {
  name     = "tf-test-kms-key-ring-name%{random_suffix}"
  location = "us"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "tf-test-kms-key-name%{random_suffix}-1"
  key_ring = google_kms_key_ring.key_ring.id
  purpose  = "ENCRYPT_DECRYPT"
}

resource "google_kms_crypto_key_iam_binding" "discoveryengine_cmek_keyuser" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
  ]
}
`, context)
}

func testAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_discovery_engine_cmek_config" "default" {
  location            = "us"
  cmek_config_id      = "tf-test-cmek-config-id%{random_suffix}"
  kms_key             = "projects/${data.google_project.project.project_id}/locations/us/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-1"
  kms_key_version     = "projects/${data.google_project.project.project_id}/locations/us/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-1/cryptoKeyVersions/1"
  set_default         = true
  single_region_keys { 
    kms_key = "projects/${data.google_project.project.project_id}/locations/us-east1/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-s1"
  }
  single_region_keys { 
    kms_key = "projects/${data.google_project.project.project_id}/locations/us-central1/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-s2"
  }
  single_region_keys { 
    kms_key = "projects/${data.google_project.project.project_id}/locations/us-west1/keyRings/tf-test-kms-key-ring-name%{random_suffix}/cryptoKeys/tf-test-kms-key-name%{random_suffix}-s3"
  }

  depends_on = [
    google_kms_crypto_key_iam_binding.discoveryengine_cmek_keyuser,
	  google_kms_crypto_key_iam_binding.discoveryengine_cmek_keyuser_s1,
    google_kms_crypto_key_iam_binding.discoveryengine_cmek_keyuser_s2,
    google_kms_crypto_key_iam_binding.discoveryengine_cmek_keyuser_s3
  ]
}

resource "google_kms_key_ring" "key_ring" {
  name     = "tf-test-kms-key-ring-name%{random_suffix}"
  location = "us"
}

resource "google_kms_crypto_key" "crypto_key" {
  name     = "tf-test-kms-key-name%{random_suffix}-1"
  key_ring = google_kms_key_ring.key_ring.id
  purpose  = "ENCRYPT_DECRYPT"
}

resource "google_kms_key_ring" "key_ring_s1" {
  name     = "tf-test-kms-key-ring-name%{random_suffix}"
  location = "us-east1"
}

resource "google_kms_crypto_key" "crypto_key_s1" {
  name     = "tf-test-kms-key-name%{random_suffix}-s1"
  key_ring = google_kms_key_ring.key_ring_s1.id
  purpose  = "ENCRYPT_DECRYPT"
}

resource "google_kms_key_ring" "key_ring_s2" {
  name     = "tf-test-kms-key-ring-name%{random_suffix}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key_s2" {
  name     = "tf-test-kms-key-name%{random_suffix}-s2"
  key_ring = google_kms_key_ring.key_ring_s2.id
  purpose  = "ENCRYPT_DECRYPT"
}

resource "google_kms_key_ring" "key_ring_s3" {
  name     = "tf-test-kms-key-ring-name%{random_suffix}"
  location = "us-west1"
}

resource "google_kms_crypto_key" "crypto_key_s3" {
  name     = "tf-test-kms-key-name%{random_suffix}-s3"
  key_ring = google_kms_key_ring.key_ring_s3.id
  purpose  = "ENCRYPT_DECRYPT"
}

resource "google_kms_crypto_key_iam_binding" "discoveryengine_cmek_keyuser" {
  crypto_key_id = google_kms_crypto_key.crypto_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "discoveryengine_cmek_keyuser_s1" {
  crypto_key_id = google_kms_crypto_key.crypto_key_s1.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "discoveryengine_cmek_keyuser_s2" {
  crypto_key_id = google_kms_crypto_key.crypto_key_s2.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
  ]
}

resource "google_kms_crypto_key_iam_binding" "discoveryengine_cmek_keyuser_s3" {
  crypto_key_id = google_kms_crypto_key.crypto_key_s3.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com",
  ]
}
`, context)
}
