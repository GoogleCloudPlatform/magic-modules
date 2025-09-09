package dataprocmetastore_test

import (
	"context"
	"fmt"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"io"
	"net/http"
	"net/url"
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
		"project":       envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetastoreServiceTags(context),
				Check: resource.TestCheckResourceAttrSet(
					"google_dataproc_metastore_service.default", "tags.%"),
			},
			{
				ResourceName:      "google_dataproc_metastore_service.default",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"service_id", "location", "labels", "terraform_labels", "tags",
				},
			},
			{
				Config:       testAccMetastoreServiceTags(context),
				ResourceName: "google_dataproc_metastore_service.default",
				Check:        checkTagBindings(context),
			},
		},
	})
}

func testAccMetastoreServiceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataproc_metastore_service" "default" {
  service_id   = "tf-test-my-service-%{random_suffix}"
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
  labels = {
    env = "test"
  }
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}

// checkTagBindings performs the API call to verify the tag binding.
func checkTagBindings(testContext map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		// Get the instance state from the state file
		rs, ok := s.RootModule().Resources["google_dataproc_metastore_service.test"]
		if !ok {
			return fmt.Errorf("not found: %s", "google_dataproc_metastore_service.test")
		}
		if rs.Primary == nil {
			return fmt.Errorf("no primary instance found for %s", "google_dataproc_metastore_service.test")
		}

		// Get the resource ID and location from the state file
		id := rs.Primary.ID
		location := rs.Primary.Attributes["location"]
		if location == "" {
			return fmt.Errorf("location not found in state file")
		}

		// Construct the full resource name
		projectID := testContext["project"].(string)
		parent := fmt.Sprintf("//metastore.googleapis.com/projects/%s/locations/%s/services/%s", projectID, location, id)

		// Build the URL for the API call
		apiURL, err := url.Parse("https://cloudresourcemanager.googleapis.com/v3/tagBindings")
		if err != nil {
			return fmt.Errorf("failed to parse API URL: %w", err)
		}
		q := apiURL.Query()
		q.Set("parent", parent)
		apiURL.RawQuery = q.Encode()

		// Create and send the GET request
		req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		// set the Authorization header here.

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to list tag bindings: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("API request failed with status code %d and body: %s", resp.StatusCode, string(body))
		}

		// Parse the response
		var result struct {
			TagBindings []struct {
				TagValue string `json:"tagValue"`
				Parent   string `json:"parent"`
				Name     string `json:"name"`
			} `json:"tagBindings"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("failed to decode API response: %w", err)
		}

		// Verify the tag binding is present
		expectedTagValue := testContext["tagValue"].(string)
		found := false
		for _, binding := range result.TagBindings {
			if strings.HasSuffix(binding.TagValue, expectedTagValue) {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("expected tag binding with value '%s' not found on resource %s", expectedTagValue, parent)
		}

		return nil
	}
}
