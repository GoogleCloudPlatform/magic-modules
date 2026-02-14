package logging_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccLoggingBucketConfigFolder_basic(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"folder_name":   "tf-test-" + acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"original_key":  acctest.BootstrapKMSKeyInLocation(t, "us-central1").CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigFolder_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_folder_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccLoggingBucketConfigFolder_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_folder_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"folder"},
			},
			{
				Config: testAccLoggingOrganizationSettings_full(context),
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_basic(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"bucket_id":       "tf-test-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 40),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_analyticsEnabled(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"bucket_id":       "tf-test-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_analyticsEnabled(context, true),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_analyticsEnabled(context, false),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigProject_cmekSettings(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"bucket_id":       "tf-test-bucket-" + acctest.RandString(t, 10),
	}

	keyRingName := fmt.Sprintf("tf-test-key-ring-%s", acctest.RandString(t, 10))
	cryptoKeyName := fmt.Sprintf("tf-test-crypto-key-%s", acctest.RandString(t, 10))
	cryptoKeyNameUpdate := fmt.Sprintf("tf-test-crypto-key-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_cmekSettings(context, keyRingName, cryptoKeyName, cryptoKeyNameUpdate),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_cmekSettingsUpdate(context, keyRingName, cryptoKeyName, cryptoKeyNameUpdate),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccLoggingBucketConfigBillingAccount_basic(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"billing_account_name": "billingAccounts/" + envvar.GetTestMasterBillingAccountFromEnv(t),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"bucket_id":            "_Default",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigBillingAccount_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_billing_account_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account"},
			},
			{
				Config: testAccLoggingBucketConfigBillingAccount_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_billing_account_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account"},
			},
		},
	})
}

func TestAccLoggingBucketConfigOrganization_basic(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"bucket_id":     "_Default",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigOrganization_basic(context, 30),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
			{
				Config: testAccLoggingBucketConfigOrganization_basic(context, 20),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
		},
	})
}

func testAccLoggingBucketConfigFolder_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
// Reset the default bucket and location settings, which may have been changed by other tests.
resource "google_logging_organization_settings" "default" {
  organization = "%{org_id}"
}

resource "google_folder" "default" {
	display_name = "%{folder_name}"
	parent       = "organizations/%{org_id}"
	deletion_protection = false
	depends_on = [google_logging_organization_settings.default]
}

resource "google_logging_folder_bucket_config" "basic" {
	folder    = google_folder.default.name
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingBucketConfigProject_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
// Reset the default bucket and location settings, which may have been changed by other tests.
resource "google_logging_organization_settings" "default" {
  organization = "%{org_id}"
}

resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
	billing_account = "%{billing_account}"
	deletion_policy = "DELETE"
}

resource "google_logging_project_bucket_config" "basic" {
	project    = google_project.default.name
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "%{bucket_id}"
}
`, context), retention, retention)
}

func testAccLoggingBucketConfigProject_analyticsEnabled(context map[string]interface{}, analytics bool) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
	billing_account = "%{billing_account}"
	deletion_policy = "DELETE"
}

// time_sleep would allow for permissions to be granted before creating log bucket
resource "time_sleep" "wait_1_minute" {
	create_duration = "1m"
  
	depends_on = [
	  google_project.default,
	]
  }

resource "google_logging_project_bucket_config" "basic" {
	project    = google_project.default.name
	location  = "global"
	enable_analytics = %t
	bucket_id = "%{bucket_id}"

	depends_on = [time_sleep.wait_1_minute]
}
`, context), analytics)
}

func testAccLoggingBucketConfigProject_locked(context map[string]interface{}, locked bool) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id = "%{project_name}"
	name       = "%{project_name}"
	org_id     = "%{org_id}"
	billing_account = "%{billing_account}"
	deletion_policy = "DELETE"
}

resource "google_logging_project_bucket_config" "fixed_locked" {
	project    = google_project.default.name
	location  = "global"
	locked = true
	bucket_id = "fixed-locked"
}

resource "google_logging_project_bucket_config" "variable_locked" {
	project    = google_project.default.name
	location  = "global"
	description = "lock status is %v" # test simultaneous update
	locked = %t
	bucket_id = "variable-locked"
}
`, context), locked, locked)
}

func testAccLoggingBucketConfigProject_preCmekSettings(context map[string]interface{}, keyRingName, cryptoKeyName, cryptoKeyNameUpdate string) string {
	return fmt.Sprintf(acctest.Nprintf(`
resource "google_project" "default" {
	project_id      = "%{project_name}"
	name            = "%{project_name}"
	org_id          = "%{org_id}"
	billing_account = "%{billing_account}"
	deletion_policy = "DELETE"
}

resource "google_project_service" "logging_service" {
	project = google_project.default.project_id
	service = "logging.googleapis.com"
}

data "google_logging_project_cmek_settings" "cmek_settings" {
	project = google_project_service.logging_service.project
}

resource "google_kms_key_ring" "keyring" {
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "key1" {
	name            = "%s"
	key_ring        = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key" "key2" {
	name            = "%s"
	key_ring        = google_kms_key_ring.keyring.id
}

resource "google_kms_crypto_key_iam_member" "crypto_key_member1" {
	crypto_key_id = google_kms_crypto_key.key1.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	member = "serviceAccount:${data.google_logging_project_cmek_settings.cmek_settings.service_account_id}"
}

resource "google_kms_crypto_key_iam_member" "crypto_key_member2" {
	crypto_key_id = google_kms_crypto_key.key2.id
	role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
	
	member = "serviceAccount:${data.google_logging_project_cmek_settings.cmek_settings.service_account_id}"
}
`, context), keyRingName, cryptoKeyName, cryptoKeyNameUpdate)
}

func testAccLoggingBucketConfigProject_cmekSettings(context map[string]interface{}, keyRingName, cryptoKeyName, cryptoKeyNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "google_logging_project_bucket_config" "basic" {
	project        = google_project.default.name
	location       = "us-central1"
	retention_days = 30
	description    = "retention test 30 days"
	bucket_id      = "%s"

	cmek_settings {
		kms_key_name = google_kms_crypto_key.key1.id
	}

	depends_on   = [google_kms_crypto_key_iam_member.crypto_key_member1]
}
`, testAccLoggingBucketConfigProject_preCmekSettings(context, keyRingName, cryptoKeyName, cryptoKeyNameUpdate), context["bucket_id"])
}

func testAccLoggingBucketConfigProject_cmekSettingsUpdate(context map[string]interface{}, keyRingName, cryptoKeyName, cryptoKeyNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "google_logging_project_bucket_config" "basic" {
	project        = google_project.default.name
	location       = "us-central1"
	retention_days = 30
	description    = "retention test 30 days"
	bucket_id      = "%s"

	cmek_settings {
		kms_key_name = google_kms_crypto_key.key2.id
	}

	depends_on   = [google_kms_crypto_key_iam_member.crypto_key_member2]
}
`, testAccLoggingBucketConfigProject_preCmekSettings(context, keyRingName, cryptoKeyName, cryptoKeyNameUpdate), context["bucket_id"])
}

func TestAccLoggingBucketConfig_CreateBuckets_withCustomId(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"random_suffix":        acctest.RandString(t, 10),
		"billing_account_name": envvar.GetTestBillingAccountFromEnv(t),
		"org_id":               envvar.GetTestOrgFromEnv(t),
		"project_name":         "tf-test-" + acctest.RandString(t, 10),
		"bucket_id":            "tf-test-bucket-" + acctest.RandString(t, 10),
	}

	configList := getLoggingBucketConfigs(context)

	for res, config := range configList {
		acctest.VcrTest(t, resource.TestCase{
			PreCheck:                 func() { acctest.AccTestPreCheck(t) },
			ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			Steps: []resource.TestStep{
				{
					Config: config,
				},
				{
					ResourceName:            fmt.Sprintf("google_logging_%s_bucket_config.basic", res),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{res},
				},
			},
		})
	}
}

func testAccLoggingBucketConfigBillingAccount_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
data "google_billing_account" "default" {
	billing_account = "%{billing_account_name}"
}

resource "google_logging_billing_account_bucket_config" "basic" {
	billing_account    = data.google_billing_account.default.billing_account
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func testAccLoggingBucketConfigOrganization_basic(context map[string]interface{}, retention int) string {
	return fmt.Sprintf(acctest.Nprintf(`
// Reset the default bucket and location settings, which may have been changed by other tests.
resource "google_logging_organization_settings" "default" {
  organization = "%{org_id}"
}

data "google_organization" "default" {
	organization = "%{org_id}"
}

resource "google_logging_organization_bucket_config" "basic" {
	organization    = data.google_organization.default.organization
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "_Default"
}
`, context), retention, retention)
}

func getLoggingBucketConfigs(context map[string]interface{}) map[string]string {
	return map[string]string{
		"project": acctest.Nprintf(`resource "google_project" "default" {
				project_id = "%{project_name}"
				name       = "%{project_name}"
				org_id     = "%{org_id}"
				billing_account = "%{billing_account_name}"
				deletion_policy = "DELETE"
			}
			
			resource "google_logging_project_bucket_config" "basic" {
				project    = google_project.default.name
				location  = "global"
				retention_days = 10
				description = "retention test 10 days"
				bucket_id = "%{bucket_id}"
			}`, context),
	}

}

func TestAccLoggingBucketConfigOrganization_indexConfigs(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"bucket_id":     "_Default",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigOrganization_indexConfigs(context, "INDEX_TYPE_STRING", "INDEX_TYPE_STRING"),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
			{
				Config: testAccLoggingBucketConfigOrganization_indexConfigs(context, "INDEX_TYPE_STRING", "INDEX_TYPE_INTEGER"),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization"},
			},
		},
	})
}

func testAccLoggingBucketConfigOrganization_indexConfigs(context map[string]interface{}, urlIndexType, statusIndexType string) string {
	return fmt.Sprintf(acctest.Nprintf(`
data "google_organization" "default" {
	organization = "%{org_id}"
}

resource "google_logging_organization_bucket_config" "basic" {
	organization    = data.google_organization.default.organization
	location  = "global"
	retention_days = 30
	description = "retention test 30 days"
	bucket_id = "_Default"

	index_configs {
		field_path 	= "jsonPayload.request.url"
		type		= "%s"
	}

	index_configs {
		field_path 	= "jsonPayload.response.status"
		type		= "%s"
	}
}
`, context), urlIndexType, statusIndexType)
}

func TestAccLoggingBucketConfigProject_indexConfigs(t *testing.T) {
	// google_logging_organization_settings is a singleton, and multiple tests mutate it.
	orgSettingsMu.Lock()
	t.Cleanup(orgSettingsMu.Unlock)

	context := map[string]interface{}{
		"project_name":    "tf-test-" + acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"bucket_id":       "tf-test-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigProject_indexConfigs(context, "INDEX_TYPE_STRING", "INDEX_TYPE_STRING"),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccLoggingBucketConfigProject_indexConfigs(context, "INDEX_TYPE_STRING", "INDEX_TYPE_INTEGER"),
			},
			{
				ResourceName:            "google_logging_project_bucket_config.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func testAccLoggingBucketConfigProject_indexConfigs(context map[string]interface{}, urlIndexType, statusIndexType string) string {
	return fmt.Sprintf(acctest.Nprintf(`
// Reset the default bucket and location settings, which may have been changed by other tests.
resource "google_logging_organization_settings" "default" {
  organization = "%{org_id}"
}

resource "google_project" "default" {
	project_id      = "%{project_name}"
	name            = "%{project_name}"
	org_id          = "%{org_id}"
	billing_account = "%{billing_account}"
	deletion_policy = "DELETE"
}

resource "google_logging_project_bucket_config" "basic" {
	project        	= google_project.default.name
	location       	= "us-east1"
	retention_days 	= 30
	description    	= "retention test 30 days"
	bucket_id      	= "%{bucket_id}"

	index_configs {
		field_path 	= "jsonPayload.request.url"
		type		= "%s"
	}

	index_configs {
		field_path 	= "jsonPayload.response.status"
		type		= "%s"
	}
}
`, context), urlIndexType, statusIndexType)
}

func TestAccLoggingBucketConfigOrganization_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "logging-bucket-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "logging-bucket-tagvalue", tagKey),
		"bucket_id":     "_Default",
	}
	retentionValue := 30

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBucketConfigOrganizationWithTags(context, retentionValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_logging_organization_bucket_config.test", "tags.%"),
					checkLoggingBucketConfigOrganizationWithTags(t),
				),
			},
			{
				ResourceName:            "google_logging_organization_bucket_config.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func checkLoggingBucketConfigOrganizationWithTags(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_organization_bucket_config" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			// 1. Get the configured tag key and value from the state.
			var configuredTagValueNamespacedName string
			for key, val := range rs.Primary.Attributes {
				if strings.HasPrefix(key, "tags.") && key != "tags.%" {
					tfTagKey := strings.TrimPrefix(key, "tags.")
					tfTagValue := val
					if tfTagValue != "" {
						configuredTagValueNamespacedName = fmt.Sprintf("%s/%s", tfTagKey, tfTagValue)
						break
					}
				}
			}

			if configuredTagValueNamespacedName == "" {
				return fmt.Errorf("could not find a configured tag value in the state for resource %s", rs.Primary.ID)
			}

			// Check if placeholders are still present.
			if strings.Contains(configuredTagValueNamespacedName, "%{") {
				return fmt.Errorf("tag namespaced name contains unsubstituted variables: %q. Ensure the context map in the test step is populated", configuredTagValueNamespacedName)
			}

			// 2. Describe the tag value using the namespaced name to get its full resource name.
			safeNamespacedName := url.QueryEscape(configuredTagValueNamespacedName)
			describeTagValueURL := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagValues/namespaced?name=%s", safeNamespacedName)

			respDescribe, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    describeTagValueURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error describing tag value using namespaced name %q: %v", configuredTagValueNamespacedName, err)
			}

			fullTagValueName, ok := respDescribe["name"].(string)
			if !ok || fullTagValueName == "" {
				return fmt.Errorf("tag value details (name) not found in response for namespaced name: %q, response: %v", configuredTagValueNamespacedName, respDescribe)
			}

			// 3. Get the tag bindings from the Logging Buckets.
			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("invalid resource ID format: %s", rs.Primary.ID)
			}
			orgID := parts[1]
			location := parts[3]
			bucket_id := parts[5]

			parentURL := fmt.Sprintf("//logging.googleapis.com/organizations/%s/locations/%s/buckets/%s", orgID, location, bucket_id)
			crmLocation := location
			if crmLocation == "global" {
				crmLocation = "us-central1"
			}
			listBindingsURL := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", crmLocation, url.QueryEscape(parentURL))

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    listBindingsURL,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("error calling TagBindings API: %v", err)
			}

			tagBindingsVal, exists := resp["tagBindings"]
			if !exists {
				tagBindingsVal = []interface{}{}
			}

			tagBindings, ok := tagBindingsVal.([]interface{})
			if !ok {
				return fmt.Errorf("'tagBindings' is not a slice in response for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			// 4. Perform the comparison.
			foundMatch := false
			for _, binding := range tagBindings {
				bindingMap, ok := binding.(map[string]interface{})
				if !ok {
					continue
				}
				if bindingMap["tagValue"] == fullTagValueName {
					foundMatch = true
					break
				}
			}

			if !foundMatch {
				return fmt.Errorf("expected tag value %s (from namespaced %q) not found in tag bindings for resource %s. Bindings: %v", fullTagValueName, configuredTagValueNamespacedName, rs.Primary.ID, tagBindings)
			}

			t.Logf("Successfully found matching tag binding for %s with tagValue %s", rs.Primary.ID, fullTagValueName)
		}

		return nil
	}
}

func testAccLoggingBucketConfigOrganizationWithTags(context map[string]interface{}, retention int) string {
	template := acctest.Nprintf(`
resource "google_logging_organization_settings" "default" {
  organization = "%{org}"
}

data "google_organization" "default" {
	organization = "%{org}"
}

resource "google_logging_organization_bucket_config" "test" {
	organization    = data.google_organization.default.organization
	location  = "global"
	retention_days = %d
	description = "retention test %d days"
	bucket_id = "%{bucket_id}"
	tags = {
	  "%{org}/%{tagKey}" = "%{tagValue}"
  }
}`, context)
	return fmt.Sprintf(template, retention, retention)
}
