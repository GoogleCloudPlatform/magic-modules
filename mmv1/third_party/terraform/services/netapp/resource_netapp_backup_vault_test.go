package netapp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetappBackupVault_NetappBackupVaultExample_update(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackupVault_NetappBackupVaultExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_backup_vault.test_backup_vault",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappBackupVault_NetappBackupVaultExample_update(context),
			},
			{
				ResourceName:            "google_netapp_backup_vault.test_backup_vault",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetappBackupVault_NetappBackupVaultExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_vault" "test_backup_vault" {
  name = "tf-test-test-backup-vault%{random_suffix}"
  location = "us-east4"
}
`, context)
}

func testAccNetappBackupVault_NetappBackupVaultExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_vault" "test_backup_vault" {
  name = "tf-test-test-backup-vault%{random_suffix}"
  location = "us-east4"
  description = "Terraform created vault"
  labels = { 
    "creator": "testuser",
	"foo": "bar",
  }
}
`, context)
}

// TestAccNetappBackupVault_Kms: Tests Backup Vault creation with KMS configuration.
func TestAccNetappBackupVault_Kms(t *testing.T) {
	location := "us-east4"
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"location":      location,
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", location, "tf-test-netapp-bv-key").CryptoKey.Name,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappBackupVault_withKms(context),
			},
			{
				ResourceName:            "google_netapp_backup_vault.test_backup_vault_kms",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetappBackupVault_withKms(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_kmsconfig" "test_kms_config" {
  name            = "tf-test-kms-config-%{random_suffix}"
  provider        = google-beta
  location        = "%{location}"
  crypto_key_name = "%{kms_key_name}"
  description     = "Test KMS config for Backup Vault"
}

resource "google_netapp_backup_vault" "test_backup_vault_kms" {
  name        = "tf-test-bv-kms-%{random_suffix}"
  provider    = google-beta
  location    = "%{location}"
  kms_config  = google_netapp_kmsconfig.test_kms_config.id
  description = "Vault with KMS"
}
	`, context)
}

func testAccCheckNetappBackupVaultDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_netapp_backup_vault" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{NetappBasePath}}projects/{{project}}/locations/{{location}}/backupVaults/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("NetappBackupVault still exists at %s", url)
			}
		}

		return nil
	}
}
