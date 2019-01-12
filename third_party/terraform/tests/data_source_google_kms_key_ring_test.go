package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleKmsKeyRing_basic(t *testing.T) {
	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check:  resource.TestMatchResourceAttr("data.google_kms_key_ring.kms_key_ring", "name", regexp.MustCompile(keyRingName)),
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
