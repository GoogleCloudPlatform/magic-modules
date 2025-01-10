package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageFolderManagementHub_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccStorageFolderManagementHub_basic(context),
			},
			{
				ResourceName:            "google_storage_folder_management_hub.folder_management_hub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageFolderManagementHub_update_with_filter(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.0.bucket_id", "random-test1"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.1.bucket_id_regex", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.2.bucket_id", "random-test2"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.3.bucket_id_regex", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.included_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.included_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_folder_management_hub.folder_management_hub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageFolderManagementHub_update_with_filter2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.0.bucket_id", "random-test1"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.1.bucket_id_regex", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.2.bucket_id", "random-test2"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.3.bucket_id_regex", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.excluded_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "filter.0.excluded_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_folder_management_hub.folder_management_hub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageFolderManagementHub_update_mode_disable(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "edition_config", "DISABLED"),
				),
			},
			{
				ResourceName:            "google_storage_folder_management_hub.folder_management_hub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageFolderManagementHub_update_mode_inherit(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_folder_management_hub.folder_management_hub", "edition_config", "INHERIT"),
				),
			},
			{
				ResourceName:            "google_storage_folder_management_hub.folder_management_hub",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccStorageFolderManagementHub_basic(context map[string]interface{}) string {
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
`, context)
}

func testAccStorageFolderManagementHub_update_with_filter(context map[string]interface{}) string {
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
  filter {
    excluded_cloud_storage_buckets{
      cloud_storage_buckets {
        bucket_id = "random-test1"
      }
      cloud_storage_buckets {
        bucket_id_regex = "random-test-*"
      }
			cloud_storage_buckets {
        bucket_id = "random-test2"
      }
      cloud_storage_buckets {
        bucket_id_regex = "random-test2-*"
      }
    }
    included_cloud_storage_locations{
      locations = ["us-east-1", "us-east-2"]
    }
  }
	depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageFolderManagementHub_update_with_filter2(context map[string]interface{}) string {
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
  filter {
    included_cloud_storage_buckets{
      cloud_storage_buckets {
        bucket_id = "random-test1"
      }
      cloud_storage_buckets {
        bucket_id_regex = "random-test-*"
      }
			cloud_storage_buckets {
        bucket_id = "random-test2"
      }
      cloud_storage_buckets {
        bucket_id_regex = "random-test2-*"
      }
    }
    excluded_cloud_storage_locations{
      locations = ["us-east-1", "us-east-2"]
    }
  }
	depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageFolderManagementHub_update_mode_disable(context map[string]interface{}) string {
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
  edition_config = "DISABLED"
	depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}

func testAccStorageFolderManagementHub_update_mode_inherit(context map[string]interface{}) string {
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
  edition_config = "INHERIT"
	depends_on = [time_sleep.wait_120_seconds]
}
`, context)
}
