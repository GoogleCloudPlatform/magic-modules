package observability_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccObservabilityFolderSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"kms_key_name":  acctest.BootstrapKMSKeyInLocation(t, "us").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccObservabilityFolderSettings_basic(context),
			},
			{
				ResourceName:            "google_observability_folder_settings.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder", "location"},
			},
			{
				Config: testAccObservabilityFolderSettings_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_observability_folder_settings.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_observability_folder_settings.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder", "location"},
			},
		},
	})
}

func testAccObservabilityFolderSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "test_folder" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

data "google_observability_folder_settings" "settings_data" {
  folder   = google_folder.test_folder.folder_id
  location = "us"
  depends_on = [google_folder.test_folder]
}

# Add a delay to allow the service account to propagate
resource "time_sleep" "wait_for_sa_propagation" {
  create_duration = "90s"
  depends_on = [data.google_observability_folder_settings.settings_data]
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_observability_folder_settings.settings_data.service_account_id}"
  depends_on = [time_sleep.wait_for_sa_propagation]
}

resource "google_observability_folder_settings" "primary" {
  location = "us"
  folder   = google_folder.test_folder.folder_id
  kms_key_name = "%{kms_key_name}"
  depends_on = [google_kms_crypto_key_iam_member.iam]
}
`, context)
}

func testAccObservabilityFolderSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "test_folder" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
  deletion_protection = false
}

data "google_observability_folder_settings" "settings_data" {
  folder   = google_folder.test_folder.folder_id
  location = "us"
  depends_on = [google_folder.test_folder]
}

# Add a delay to allow the service account to propagate
resource "time_sleep" "wait_for_sa_propagation" {
  create_duration = "90s"
  depends_on = [data.google_observability_folder_settings.settings_data]
}

resource "google_kms_crypto_key_iam_member" "iam" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:${data.google_observability_folder_settings.settings_data.service_account_id}"
  depends_on = [time_sleep.wait_for_sa_propagation]
}

resource "google_observability_folder_settings" "primary" {
  location = "us"
  folder   = google_folder.test_folder.folder_id
  kms_key_name = "" # Unset KMS key
  depends_on = [google_kms_crypto_key_iam_member.iam]
}
`, context)
}

func TestAccObservabilityFolderSettings_globalUpdate(t *testing.T) {
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
				Config: testAccObservabilityFolderSettings_global(context, "us"),
			},
			{
				ResourceName:            "google_observability_folder_settings.primary_global",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder", "location"},
			},
			{
				Config: testAccObservabilityFolderSettings_globalUpdate(context, "eu"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_observability_folder_settings.primary_global", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_observability_folder_settings.primary_global",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder", "location"},
			},
		},
	})
}

func testAccObservabilityFolderSettings_global(context map[string]interface{}, defaultLocation string) string {
	return acctest.Nprintf(fmt.Sprintf(`
resource "google_folder" "test_folder" {
  display_name = "tf-test-%%{random_suffix}"
  parent       = "organizations/%%{org_id}"
  deletion_protection = false
}

data "google_observability_folder_settings" "settings_data" {
  folder   = google_folder.test_folder.folder_id
  location = "global"
  depends_on = [google_folder.test_folder]
}

resource "google_observability_folder_settings" "primary_global" {
  location                 = "global"
  folder                   = google_folder.test_folder.folder_id
  default_storage_location = "%s"
  depends_on = [data.google_observability_folder_settings.settings_data]
}
`, defaultLocation), context)
}

func testAccObservabilityFolderSettings_globalUpdate(context map[string]interface{}, defaultLocation string) string {
	return acctest.Nprintf(fmt.Sprintf(`
resource "google_folder" "test_folder" {
  display_name = "tf-test-%%{random_suffix}"
  parent       = "organizations/%%{org_id}"
  deletion_protection = false
}

data "google_observability_folder_settings" "settings_data" {
  folder   = google_folder.test_folder.folder_id
  location = "global"
  depends_on = [google_folder.test_folder]
}

resource "google_observability_folder_settings" "primary_global" {
  location                 = "global"
  folder                   = google_folder.test_folder.folder_id
  default_storage_location = "%s"
  depends_on = [data.google_observability_folder_settings.settings_data]
}
`, defaultLocation), context)
}
