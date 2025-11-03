package backupdr_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
	"time"
)

func TestAccBackupDRBackupVault_fullUpdate(t *testing.T) {
	// Uses time.Now
	acctest.SkipIfVcr(t)

	t.Parallel()

	timeNow := time.Now().UTC()
	referenceTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC)

	context := map[string]interface{}{
		"project":        envvar.GetTestProjectFromEnv(),
		"effective_time": referenceTime.Add(24 * time.Hour).Format(time.RFC3339),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupDRBackupVault_fullCreate(context),
			},
			{
				ResourceName:            "google_backup_dr_backup_vault.backup-vault-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_missing", "annotations", "backup_vault_id", "force_delete", "force_update", "ignore_backup_plan_references", "ignore_inactive_datasources", "backup_retention_inheritance", "access_restriction", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccBackupDRBackupVault_fullUpdate(context),
			},
			{
				ResourceName:            "google_backup_dr_backup_vault.backup-vault-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_missing", "annotations", "backup_vault_id", "force_delete", "force_update", "ignore_backup_plan_references", "ignore_inactive_datasources", "backup_retention_inheritance", "access_restriction", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func TestAccBackupDRBackupVault_createUsingCMEK(t *testing.T) {
	// Uses time.Now
	acctest.SkipIfVcr(t)

	t.Parallel()

	timeNow := time.Now().UTC()
	referenceTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.UTC)

	context := map[string]interface{}{
		"project":        envvar.GetTestProjectFromEnv(),
		"effective_time": referenceTime.Add(24 * time.Hour).Format(time.RFC3339),
		"random_suffix":  acctest.RandString(t, 10),
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-bv-cmek").CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupDRBackupVault_createUsingCMEK(context),
			},
			{
				ResourceName:            "google_backup_dr_backup_vault.backup-vault-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_missing", "annotations", "backup_vault_id", "force_delete", "force_update", "ignore_backup_plan_references", "ignore_inactive_datasources", "backup_retention_inheritance", "access_restriction", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccBackupDRBackupVault_fullCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_backup_dr_backup_vault" "backup-vault-test" {
  location = "us-central1"
  backup_vault_id    = "tf-test-backup-vault-test%{random_suffix}"
  description = "This is a backup vault built by Terraform."
  backup_minimum_enforced_retention_duration = "100000s"
  effective_time = "%{effective_time}" 
  labels = {
    foo = "bar"
	bar = "baz"
  }
  annotations = {
    annotations1 = "bar"
	annotations2 = "baz"
  }
  force_update = "true"
  ignore_inactive_datasources = "true"
  access_restriction = "WITHIN_ORGANIZATION"
  backup_retention_inheritance = "INHERIT_VAULT_RETENTION"
  ignore_backup_plan_references = "true"
  allow_missing = "true"
}
`, context)
}

func testAccBackupDRBackupVault_fullUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_backup_dr_backup_vault" "backup-vault-test" {
  location = "us-central1"
  backup_vault_id    = "tf-test-backup-vault-test%{random_suffix}"
  description = "This is a second backup vault built by Terraform."
  backup_minimum_enforced_retention_duration = "200000s"
  effective_time = "%{effective_time}" 
  labels = {
	foo = "bar1"
	bar = "baz1"
  }
  annotations = {
    annotations1 = "bar1"
	annotations2 = "baz1"
  }
  force_update = "true"
  access_restriction = "WITHIN_ORGANIZATION"
  backup_retention_inheritance = "INHERIT_VAULT_RETENTION"
  ignore_inactive_datasources = "true"
  ignore_backup_plan_references = "true"
  allow_missing = "true"
}
`, context)
}

func testAccBackupDRBackupVault_createUsingCMEK(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_backup_dr_backup_vault" "backup-vault-test-cmek" {
  provider      = google-beta
  location = "us-central1"
  backup_vault_id    = "tf-test-backup-vault-test%{random_suffix}"
  description = "This is a second backup vault built by Terraform."
  backup_minimum_enforced_retention_duration = "200000s"
  effective_time = "%{effective_time}" 
  labels = {
	foo = "bar1"
	bar = "baz1"
  }
  annotations = {
    annotations1 = "bar1"
	annotations2 = "baz1"
  }
  encryption_config {
    kms_key_name = "%{{kms_key_name}}"
  }
  force_update = "true"
  access_restriction = "WITHIN_ORGANIZATION"
  backup_retention_inheritance = "INHERIT_VAULT_RETENTION"
  ignore_inactive_datasources = "true"
  ignore_backup_plan_references = "true"
  allow_missing = "true"
  # Ensure IAM binding is created before the vault attempts to use the key
  depends_on = [
    google_kms_crypto_key_iam_member.backupdr_sa_crypto_access
  ]
}

# Grant the BackupDR P4SA permission to use the key
resource "google_kms_crypto_key_iam_member" "backupdr_sa_crypto_access" {
  provider      = google-beta
  crypto_key_id = "%{{kms_key_name}}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.test_project.number}@gcp-sa-backupdr.iam.gserviceaccount.com"
}
`, context)
}
