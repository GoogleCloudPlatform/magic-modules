package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventarcChannel_cryptoKeyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":         envvar.GetTestRegionFromEnv(),
		"project_name":   envvar.GetTestProjectFromEnv(),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcChannel_setCryptoKey(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcChannel_cryptoKeyUpdate(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcChannel_setCryptoKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_key_ring" "test_key_ring" {
	name     = "tf-keyring%{random_suffix}"
	location = "us-central1"
}

resource "google_kms_crypto_key" "key1" {
	name     = "tf-key1%{random_suffix}"
	key_ring = google_kms_key_ring.test_key_ring.id
}

resource "google_kms_crypto_key_iam_member" "key1_member" {
	crypto_key_id = google_kms_crypto_key.key1.id
	role      = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

	member = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_channel" "primary" {
	location = "%{region}"
	name     = "tf-test-name%{random_suffix}"
	crypto_key_name =  google_kms_crypto_key.key1.id
	third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
	depends_on = [google_kms_crypto_key_iam_member.key1_member]
}
`, context)
}

func testAccEventarcChannel_cryptoKeyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_key_ring" "test_key_ring" {
	name     = "tf-keyring%{random_suffix}"
	location = "us-central1"
}
	
resource "google_kms_crypto_key" "key2" {
	name     = "tf-key2%{random_suffix}"
	key_ring = google_kms_key_ring.test_key_ring.id
}

resource "google_kms_crypto_key_iam_member" "key2_member" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role      = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	member = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_channel" "primary" {
	location = "%{region}"
	name     = "tf-test-name%{random_suffix}"
	crypto_key_name = google_kms_crypto_key.key2.id
	third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
	depends_on = [google_kms_crypto_key_iam_member.key2_member]
}
`, context)
}

func TestAccEventarcChannel_LongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcChannel_LongForm(context),
			},
			{
				ResourceName:            "google_eventarc_channel.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "project"},
			},
		},
	})
}

func testAccEventarcChannel_LongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_channel" "primary" {
	location = "long/form/%{region}"
	project  = "projects/%{project_name}"
	name     = "projects/%{project_name}/locations/%{region}/channels/tf-test-name%{random_suffix}"
	third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
}
`, context)
}
