package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingFolderSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSettings_full(context),
			},
			{
				ResourceName:            "google_logging_folder_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccLoggingFolderSettings_update(context),
			},
			{
				ResourceName:            "google_logging_folder_settings.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
		},
	})
}

func testAccLoggingFolderSettings_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_folder_settings" "example" {
  disable_default_sink = true
  folder               = google_folder.my_folder.folder_id
  kms_key_name         = google_kms_crypto_key.key.id
  storage_location     = "us-central1"
  depends_on           = [ google_kms_crypto_key_iam_member.iam ]
}

resource "google_kms_key_ring" "keyring" {
  name     = "tf-test-keyring-%{random_suffix}"
  location = "us-central1"
}

resource "google_kms_crypto_key" "key" {
  name            = "tf-test-k-%{random_suffix}"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"
}

resource "google_folder" "my_folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

data "google_logging_folder_settings" "settings" {
  folder = google_folder.my_folder.folder_id
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.key.id
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:${data.google_logging_folder_settings.settings.kms_service_account_id}"
}
`, context)
}

func testAccLoggingFolderSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_folder_settings" "example" {
  disable_default_sink = false
  folder               = google_folder.my_folder.folder_id
  kms_key_name         = google_kms_crypto_key.key.id
  storage_location     = "us-east1"
  depends_on           = [ google_kms_crypto_key_iam_member.iam ]
}

resource "google_kms_key_ring" "keyring" {
  name     = "tf-test-keyring-%{random_suffix}"
  location = "us-east1"
}

resource "google_kms_crypto_key" "key" {
  name            = "tf-test-k-%{random_suffix}"
  key_ring        = google_kms_key_ring.keyring.id
  rotation_period = "100000s"
}

resource "google_folder" "my_folder" {
  display_name = "tf-test-folder-%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

data "google_logging_folder_settings" "settings" {
  folder = google_folder.my_folder.folder_id
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = google_kms_crypto_key.key.id
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:${data.google_logging_folder_settings.settings.kms_service_account_id}"
}
`, context)
}
