package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleKmsKeyRing_basic(t *testing.T) {
	kms := BootstrapKMSKey(t)
	projectId := getTestProjectFromEnv()
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, kms.KeyRing.Name),
				Check:  resource.TestMatchResourceAttr("data.google_kms_key_ring.kms_key_ring", "name", regexp.MustCompile(kms.KeyRing.Name)),
			},
		},
	})
}

/*
	This test should run in its own project, because keys and key rings are not deletable
*/
func testAccDataSourceGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
	name			= "%s"
	project_id		= "%s"
	org_id			= "%s"
	billing_account	= "%s"
}

resource "google_project_services" "acceptance" {
	project  = "${google_project.acceptance.project_id}"
	services = [
		"cloudkms.googleapis.com"
	]
}

resource "google_kms_key_ring" "key_ring" {
	project  = "${google_project_services.acceptance.project}"
	name     = "%s"
	location = "us-central1"
}

data "google_kms_key_ring" "kms_key_ring" {
	name     = "%s"
	location = "us-central1"
}
	`, projectId, projectId, projectOrg, projectBillingAccount, keyRingName, keyRingName)
}
