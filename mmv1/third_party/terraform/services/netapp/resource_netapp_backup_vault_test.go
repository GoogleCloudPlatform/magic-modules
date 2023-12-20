package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappbackupVault_netappBackupVaultExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappbackupVault_netappBackupVaultExample_basic(context),
			},
			{
				ResourceName:            "google_netapp_backup_vault.test_backup_vault",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappbackupVault_netappBackupVaultExample_update(context),
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

func testAccNetappbackupVault_netappBackupVaultExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_vault" "test_backup_vault" {
  name = "tf-test-test-backup-vault%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func testAccNetappbackupVault_netappBackupVaultExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_netapp_backup_vault" "test_backup_vault" {
  name = "tf-test-test-backup-vault%{random_suffix}"
  location = "us-central1"
  description = "Terraform created vault"
  labels = { 
    "creator": "testuser"
  }
}
`, context)
}
