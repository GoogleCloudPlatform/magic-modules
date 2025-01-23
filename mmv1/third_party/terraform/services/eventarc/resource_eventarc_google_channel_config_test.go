package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventarcGoogleChannelConfig_cryptoKeyUpdate(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"project_name":   envvar.GetTestProjectFromEnv(),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"region":         region,
		"random_suffix":  acctest.RandString(t, 10),
		"key1":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-google-channel-config-key1").CryptoKey.Name,
		"key2":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-google-channel-config-key2").CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcGoogleChannelConfig_setCryptoKey(context),
			},
			{
				ResourceName:      "google_eventarc_google_channel_config.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcGoogleChannelConfig_cryptoKeyUpdate(context),
			},
			{
				ResourceName:      "google_eventarc_google_channel_config.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcGoogleChannelConfig_deleteCryptoKey(context),
			},
		},
	})
}

func testAccEventarcGoogleChannelConfig_setCryptoKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key1_member" {
  crypto_key_id = "%{key1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_google_channel_config" "primary" {
  location        = "%{region}"
  name            = "projects/%{project_name}/locations/%{region}/googleChannelConfig"
  crypto_key_name = "%{key1}"
  depends_on      = [google_kms_crypto_key_iam_member.key1_member]
}
`, context)
}

func testAccEventarcGoogleChannelConfig_cryptoKeyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key1_member" {
  crypto_key_id = "%{key1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "key2_member" {
  crypto_key_id = "%{key2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_google_channel_config" "primary" {
  location        = "%{region}"
  name            = "projects/%{project_name}/locations/%{region}/googleChannelConfig"
  crypto_key_name = "%{key2}"
  depends_on      = [google_kms_crypto_key_iam_member.key1_member, google_kms_crypto_key_iam_member.key2_member]
}
`, context)
}

func testAccEventarcGoogleChannelConfig_deleteCryptoKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key1_member" {
  crypto_key_id = "%{key1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "key2_member" {
  crypto_key_id = "%{key2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_google_channel_config" "primary" {
  location        = "%{region}"
  name            = "projects/%{project_name}/locations/%{region}/googleChannelConfig"
  crypto_key_name = ""
  depends_on      = [google_kms_crypto_key_iam_member.key1_member, google_kms_crypto_key_iam_member.key2_member]
}
`, context)
}
