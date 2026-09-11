package secretmanagerregional_test

import (
	"maps"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/secretmanagerregional"
)

func TestAccSecretManagerRegionalRegionalSecretVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalRegionalSecretVersion_basic(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_version.secret-version-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecretManagerRegionalRegionalSecretVersion_disable(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_version.secret-version-basic",
				ImportState:       true,
				ImportStateVerify: true,
				// at this point the secret data is disabled and so reading the data on import will
				// give an empty string
				ImportStateVerifyIgnore: []string{"secret_data"},
			},
			{
				Config: testAccSecretManagerRegionalRegionalSecretVersion_basic(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_version.secret-version-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecretVersion_cmekOutputOnly(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":  kms.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-secret-manager-managed-central-key5").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalRegionalSecretVersion_cmekOutputOnly(context),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_version.secret-version-cmek",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecretVersion_writeOnlyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalRegionalSecretVersion_writeOnly(context, "my-tf-test-secret", 1),
			},
			{
				ResourceName:      "google_secret_manager_regional_secret_version.secret-version-write-only",
				ImportState:       true,
				ImportStateVerify: true,
				// write-only data is never stored in state, and its version is client-side only,
				// so neither can be read back on import
				ImportStateVerifyIgnore: []string{"secret_data_wo", "secret_data_wo_version"},
			},
			{
				// bumping the version replaces the secret version with one holding the new data
				Config: testAccSecretManagerRegionalRegionalSecretVersion_writeOnly(context, "my-tf-test-secret-updated", 2),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret_version.secret-version-write-only",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_data_wo", "secret_data_wo_version"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecretVersion_neitherSecretDataSet(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccSecretManagerRegionalRegionalSecretVersion_noSecretData(context),
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
		},
	})
}

func testAccSecretManagerRegionalRegionalSecretVersion_writeOnly(context map[string]interface{}, data string, version int) string {
	context = maps.Clone(context)
	context["secret_data_wo"] = data
	context["secret_data_wo_version"] = version

	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-write-only" {
  secret = google_secret_manager_regional_secret.secret-basic.name
  secret_data_wo = "%{secret_data_wo}%{random_suffix}"
  secret_data_wo_version = %{secret_data_wo_version}
}
`, context)
}

func testAccSecretManagerRegionalRegionalSecretVersion_noSecretData(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%{random_suffix}"
  location = "us-central1"
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.name
  secret_data_wo_version = 3
}
`, context)
}

func testAccSecretManagerRegionalRegionalSecretVersion_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%{random_suffix}"
  location = "us-central1"
  labels = {
    label = "my-label"
  }
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.name
  secret_data = "my-tf-test-secret%{random_suffix}"
  enabled = true
}
`, context)
}

func testAccSecretManagerRegionalRegionalSecretVersion_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "secret-basic" {
  secret_id = "tf-test-secret-version-%{random_suffix}"
  location = "us-central1"
  labels = {
    label = "my-label"
  }
}

resource "google_secret_manager_regional_secret_version" "secret-version-basic" {
  secret = google_secret_manager_regional_secret.secret-basic.name
  secret_data = "my-tf-test-secret%{random_suffix}"
  enabled = false
}
`, context)
}

func testAccSecretManagerRegionalRegionalSecretVersion_cmekOutputOnly(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project-ds" {}

resource "google_kms_crypto_key_iam_member" "kms-secret-binding-reg-sec-ver" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project-ds.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-reg-sec-ver" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  customer_managed_encryption {
    kms_key_name = "%{kms_key_name}"
  }

  depends_on = [ google_kms_crypto_key_iam_member.kms-secret-binding-reg-sec-ver ]
}

resource "google_secret_manager_regional_secret_version" "secret-version-cmek" {
  secret = google_secret_manager_regional_secret.regional-secret-reg-sec-ver.name
  secret_data = "my-tf-test-secret%{random_suffix}"
}
`, context)
}
