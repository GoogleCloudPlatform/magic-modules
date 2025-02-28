// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleStorageFolderManagementHub_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
		"org_id":        envvar.GetTestOrgFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleStorageFolderManagementHub_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_storage_folder_management_hub.folder_management_hub", "google_storage_folder_management_hub.folder_management_hub"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleStorageFolderManagementHub_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_folder" "folder" {
  parent       = "organizations/%{org_id}"
  display_name = "tf-test-folder-name%{random_suffix}"
	deletion_protection=false
}
resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_folder.folder]
  create_duration = "120s"
}
resource "google_storage_folder_management_hub" "folder_management_hub" {
  name = google_folder.folder.folder_id
  edition_config = "STANDARD"
	depends_on = [time_sleep.wait_120_seconds]
}

data "google_storage_folder_management_hub" "folder_management_hub" {
  name = google_storage_folder_management_hub.folder_management_hub.name
}
`, context)
}
