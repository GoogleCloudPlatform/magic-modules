package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleKmsCryptoKey_basic(t *testing.T) {
	kms := BootstrapKMSKey(t)
	projectId := getTestProjectFromEnv()
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, kms.KeyRing.Name, kms.CryptoKey.Name),
				Check:  resource.TestMatchResourceAttr("data.google_kms_crypto_key", "name", regexp.MustCompile(kms.CryptoKey.Name)),
			},
		},
	})
}

/*
	This test should run in its own project, because KMS key rings and crypto keys are not deletable
*/
func testAccDataSourceGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name            = "%s"
	project_id      = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_services" "acceptance" {
	project = "${google_project.acceptance.project_id}"

	services = [
	  "cloudkms.googleapis.com",
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "crypto_key" {
	name            = "%s"
	key_ring        = "${google_kms_key_ring.key_ring.self_link}"
	rotation_period = "1000000s"
	version_template {
		algorithm =        "GOOGLE_SYMMETRIC_ENCRYPTION"
		protection_level = "SOFTWARE"
	}
}

data "google_kms_crypto_key" "kms_crypto_key" {
	name     = "%s"
	key_ring = "${google_kms_key_ring.key_ring.self_link}"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName, cryptoKeyName)
}
