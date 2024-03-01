package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccKmsKeyHandle_basic(t *testing.T) {
	projectId := envvar.GetTestProjectFromEnv()
	projectOrg := envvar.GetTestOrgFromEnv(t)
	projectBillingAccount := envvar.GetTestBillingAccountFromEnv(t)
	keyHandleName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleKmsKeyHandleWasRemovedFromState("google_kms_key_handle.key_handle"),
		Steps: []resource.TestStep{
			{
				Config: testGoogleKmsKeyHandle_basic(projectId, keyHandleName),
			},
			{
				ResourceName:      "google_kms_key_handle.key_handle",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleKmsKeyHandle_removed(projectId, projectOrg, projectBillingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleKmsKeyHandleWasRemovedFromState("google_kms_key_handle.key_handle"),
				),
			},
		},
	})
}

// KMS KeyHandles cannot be deleted. This ensures that the KeyHandle resource was removed from state,
// even though the server-side resource was not removed.
func testAccCheckGoogleKmsKeyHandleWasRemovedFromState(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]

		if ok {
			return fmt.Errorf("Resource was not removed from state: %s", resourceName)
		}

		return nil
	}
}

// This test runs in its own project, otherwise the test project would start to get filled
// with undeletable resources
func testGoogleKmsKeyHandle_basic(projectId, keyHandleName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_kms_autokey_config" "autokey_config" {
  folder      = "343188819919"
  key_project = "projects/${data.google_project.project.number}"
  enabled     = true
}

resource "google_kms_key_handle" "key_handle" {
  name     = "%s"
  location = "global"
  resource_type_selector = "compute.googleapis.com/Disk"
}
`, projectId, keyHandleName)
}

func testGoogleKmsKeyHandle_removed(projectId, projectOrg, projectBillingAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  name            = "%s"
  project_id      = "%s"
  folder_id       = "343188819919"
  billing_account = "%s"
}
`, projectId, projectId, projectBillingAccount)
}
