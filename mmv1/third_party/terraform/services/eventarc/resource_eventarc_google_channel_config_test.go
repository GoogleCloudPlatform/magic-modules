package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// We make sure not to run tests in parallel, since only one GoogleChannelConfig per location is supported.
func TestAccEventarcGoogleChannelConfig(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic":           testAccEventarcGoogleChannelConfig_basic,
		"longForm":        testAccEventarcGoogleChannelConfig_longForm,
		"cryptoKey":       testAccEventarcGoogleChannelConfig_cryptoKey,
		"cryptoKeyUpdate": testAccEventarcGoogleChannelConfig_cryptoKeyUpdate,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccEventarcGoogleChannelConfig_basic(t *testing.T) {
	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcGoogleChannelConfig_basicCfg(context),
			},
			{
				ResourceName:            "google_eventarc_google_channel_config.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccEventarcGoogleChannelConfig_basicCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_google_channel_config" "primary" {
  location = "%{region}"
  name     = "googleChannelConfig"
}
`, context)
}

func testAccEventarcGoogleChannelConfig_longForm(t *testing.T) {
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
				Config: testAccEventarcGoogleChannelConfig_longFormCfg(context),
			},
			{
				ResourceName:            "google_eventarc_google_channel_config.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "project"},
			},
		},
	})
}

func testAccEventarcGoogleChannelConfig_longFormCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_google_channel_config" "primary" {
  project  = "projects/%{project_name}"
  location = "long/form/%{region}"
  name     = "projects/%{project_name}/locations/%{region}/googleChannelConfig"
}
`, context)
}

func testAccEventarcGoogleChannelConfig_cryptoKey(t *testing.T) {
	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"region":         region,
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"key_name":       acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-google-channel-config-key").CryptoKey.Name,
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcGoogleChannelConfig_cryptoKeyCfg(context),
			},
			{
				ResourceName:            "google_eventarc_google_channel_config.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccEventarcGoogleChannelConfig_cryptoKeyCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key_member" {
  crypto_key_id = "%{key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_google_channel_config" "primary" {
  location        = "%{region}"
  name            = "googleChannelConfig"
  crypto_key_name = "%{key_name}"
  depends_on      = [google_kms_crypto_key_iam_member.key_member]
}
`, context)
}

func testAccEventarcGoogleChannelConfig_cryptoKeyUpdate(t *testing.T) {
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
				Config: testAccEventarcGoogleChannelConfig_setCryptoKeyCfg(context),
			},
			{
				ResourceName:      "google_eventarc_google_channel_config.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcGoogleChannelConfig_cryptoKeyUpdateCfg(context),
			},
			{
				ResourceName:      "google_eventarc_google_channel_config.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcGoogleChannelConfig_deleteCryptoKeyCfg(context),
			},
		},
	})
}

func testAccEventarcGoogleChannelConfig_setCryptoKeyCfg(context map[string]interface{}) string {
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

func testAccEventarcGoogleChannelConfig_cryptoKeyUpdateCfg(context map[string]interface{}) string {
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

func testAccEventarcGoogleChannelConfig_deleteCryptoKeyCfg(context map[string]interface{}) string {
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
