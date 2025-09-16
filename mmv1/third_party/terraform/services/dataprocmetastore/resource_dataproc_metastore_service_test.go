package dataprocmetastore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
)

func TestAccDataprocMetastoreService_updateAndImport(t *testing.T) {
	t.Parallel()

	name := "tf-test-metastore-" + acctest.RandString(t, 10)
	tier := [2]string{"DEVELOPER", "ENTERPRISE"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[0]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataprocMetastoreService_updateAndImport(name, tier[1]),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.my_metastore",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataprocMetastoreService_updateAndImport(name, tier string) string {
	return fmt.Sprintf(`
resource "google_dataproc_metastore_service" "my_metastore" {
	service_id = "%s"
	location   = "us-central1"
	tier       = "%s"

	hive_metastore_config {
		version = "2.3.6"
	}
}
`, name, tier)
}

func TestAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExample(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location", "labels", "terraform_labels"},
			},
			{
				Config: testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(context),
			},
		},
	})
}

func TestAccDataprocMetastoreService_PrivateServiceConnect(t *testing.T) {
	t.Skip("Skipping due to https://github.com/hashicorp/terraform-provider-google/issues/13710")
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataprocMetastoreService_PrivateServiceConnect(context),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location"},
			},
		},
	})
}

func testAccDataprocMetastoreService_PrivateServiceConnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
// Use data source instead of creating a subnetwork due to a bug on API side.
// With the bug, the new created subnetwork cannot be deleted when deleting the dataproc metastore service.
data "google_compute_subnetwork" "subnet" {
  name   = "default"
  region = "us-central1"
}

resource "google_dataproc_metastore_service" "default" {
  service_id = "tf-test-metastore-srv%{random_suffix}"
  location   = "us-central1"

  hive_metastore_config {
    version = "3.1.2"
  }

  network_config {
    consumers {
      subnetwork = data.google_compute_subnetwork.subnet.id
    }
  }
}
`, context)
}

func testAccDataprocMetastoreService_dataprocMetastoreServiceScheduledBackupExampleUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "backup" {
  service_id = "tf-test-backup%{random_suffix}"
  location   = "us-central1"
  port       = 9080
  tier       = "DEVELOPER"

  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }

  hive_metastore_config {
    version = "2.3.6"
  }

  scheduled_backup {
    enabled         = true
    cron_schedule   = "0 0 * * 0"
    time_zone       = "America/Los_Angeles"
    backup_location = "gs://${google_storage_bucket.bucket.name}"
  }

  labels = {
    env = "test"
  }
}

resource "google_storage_bucket" "bucket" {
  name     = "tf-test-backup%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func TestAccMetastoreService_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "metastore-service-tagkey", map[string]interface{}{})
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "metastore-service-tagvalue", tagKey),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataprocMetastoreServiceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreServiceTags(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_dataproc_metastore_service.default", "tags.%"),
					testAccCheckMetastoreServiceHasTagBindings(t, map[string]bool{
						context["tagValue"].(string): true,
					}),
				),
			},
			{
				ResourceName:            "google_dataproc_metastore_service.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_id", "location", "labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCheckMetastoreServiceHasTagBindings(t *testing.T, expectedTags map[string]bool) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_metastore_service" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			parts := strings.Split(rs.Primary.ID, "/")
			if len(parts) != 6 {
				return fmt.Errorf("Invalid resource ID format: %s", rs.Primary.ID)
			}
			project := parts[1]
			location := parts[3]
			service_id := parts[5]

			parentURL := fmt.Sprintf("//metastore.googleapis.com/projects/%s/locations/%s/services/%s", project, location, service_id)
			url := fmt.Sprintf("https://%s-cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", location, parentURL)
			fmt.Printf("Checking tagBindings for resource: %s at URL: %s\n", rs.Primary.ID, url)

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})

			if err != nil {
				return fmt.Errorf("Error calling TagBindings API: %v", err)
			}

			tagBindingsVal, exists := resp["tagBindings"]
			if !exists {
				return fmt.Errorf("Key 'tagBindings' not found in response for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			tagBindings, ok := tagBindingsVal.([]interface{})
			if !ok {
				return fmt.Errorf("'tagBindings' is not a slice in response for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			// Check if any tag bindings were found.
			if len(tagBindings) == 0 {
				return fmt.Errorf("No tag bindings found for resource %s. Response: %v", rs.Primary.ID, resp)
			}

			foundTags := make(map[string]bool)
		for _, tb := range tagBindings {
			tbMap, ok := tb.(map[string]interface{})
			if !ok {
				continue
			}
			if tag, ok := tbMap["tagValue"].(string); ok {
				foundTags[tag] = true
			}
		}

			for expectedTagValue := range expectedTags {
			if !foundTags[expectedTagValue] && !foundTags["tagValues/"+expectedTagValue] {
				return fmt.Errorf("Expected tag binding %s not found in resource %s. Got %v", expectedTagValue, rs.Primary.ID, foundTags)
			}
		}
		
		fmt.Printf("Successfully found %v tag bindings for %s\n", expectedTags, rs.Primary.ID)
		}

		return nil
	}
}

func testAccMetastoreServiceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`resource "google_dataproc_metastore_service" "default" {
  service_id   = "tf-test-my-service-%{random_suffix}"
  location   = "us-east1"
  port       = 9080
  tier       = "DEVELOPER"
  maintenance_window {
    hour_of_day = 2
    day_of_week = "SUNDAY"
  }
  hive_metastore_config {
    version = "2.3.6"
  }
  labels = {
    env = "test"
  }
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}`, context)
}
