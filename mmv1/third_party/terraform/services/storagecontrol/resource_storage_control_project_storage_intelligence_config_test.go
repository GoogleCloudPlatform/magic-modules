package storagecontrol_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccStorageControlProjectIntelligenceConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageControlProjectIntelligenceConfig_basic(context),
			},
			{
				ResourceName:            "google_storage_control_project_intelligence_config.project_storage_intelligence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlProjectIntelligenceConfig_update_with_filter(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.0.bucket_id", "random-test1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.1.bucket_id_regex", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.2.bucket_id", "random-test2"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.excluded_cloud_storage_buckets.0.cloud_storage_buckets.3.bucket_id_regex", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.included_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.included_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_control_project_intelligence_config.project_storage_intelligence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlProjectIntelligenceConfig_update_with_filter2(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.0.bucket_id", "random-test1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.1.bucket_id_regex", "random-test-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.2.bucket_id", "random-test2"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.included_cloud_storage_buckets.0.cloud_storage_buckets.3.bucket_id_regex", "random-test2-*"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.excluded_cloud_storage_locations.0.locations.0", "us-east-1"),
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "filter.0.excluded_cloud_storage_locations.0.locations.1", "us-east-2"),
				),
			},
			{
				ResourceName:            "google_storage_control_project_intelligence_config.project_storage_intelligence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlProjectIntelligenceConfig_update_mode_disable(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "edition_config", "DISABLED"),
				),
			},
			{
				ResourceName:            "google_storage_control_project_intelligence_config.project_storage_intelligence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccStorageControlProjectIntelligenceConfig_update_mode_inherit(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_storage_control_project_intelligence_config.project_storage_intelligence", "edition_config", "INHERIT"),
				),
			},
			{
				ResourceName:            "google_storage_control_project_intelligence_config.project_storage_intelligence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccStorageControlProjectIntelligenceConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = "%{project}"
  edition_config = "STANDARD"
}
`, context)
}

func testAccStorageControlProjectIntelligenceConfig_update_with_filter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = "%{project}"
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
}
`, context)
}

func testAccStorageControlProjectIntelligenceConfig_update_with_filter2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = "%{project}"
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
}
`, context)
}

func testAccStorageControlProjectIntelligenceConfig_update_mode_disable(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = "%{project}"
  edition_config = "DISABLED"
}
`, context)
}

func testAccStorageControlProjectIntelligenceConfig_update_mode_inherit(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_control_project_intelligence_config" "project_storage_intelligence" {
  name = "%{project}"
  edition_config = "INHERIT"
}
`, context)
}
