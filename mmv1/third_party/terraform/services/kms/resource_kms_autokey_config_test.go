package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccKmsAutokeyConfig_basic(t *testing.T) {
	projectId := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	projectOrg := envvar.GetTestOrgFromEnv(t)
	projectBillingAccount := envvar.GetTestBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsAutokeyConfig_basic(projectId, projectOrg, projectBillingAccount),
			},
			{
				Config: testGoogleKmsAutokeyConfig_update(projectId, projectOrg, projectBillingAccount),
			},
			{
				ResourceName:      "google_kms_autokey_config.autokey_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGoogleKmsAutokeyConfig_basic(projectId, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_folder" "folder" {
  display_name = "my-folder"
  parent       = "organizations/%s"
}

resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  folder_id       = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_autokey_config" "autokey_config" {
  folder      = google_folder.folder.folder_id
  key_project = google_project.acceptance.project_id
}
`, projectOrg, projectId, projectId, projectOrg, projectBillingAccount)
}

func testGoogleKmsAutokeyConfig_update(projectId, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_folder" "folder" {
  display_name = "my-folder"
  parent       = "organizations/%s"
}

resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  folder_id       = "%s"
  billing_account = "%s"
}

resource "google_project_service" "acceptance" {
  project = google_project.acceptance.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_kms_autokey_config" "autokey_config" {
  folder      = google_folder.folder.folder_id
  key_project = ""
}
`, projectOrg, projectId, projectId, projectOrg, projectBillingAccount)
}
