package dataprocmetastore_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "google.golang.org/api/transport/http"

	"io"
	"net/http"
	"encoding/json"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/test/acctest"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/test/envvar"
)

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
				Check: resource.TestCheckFunc(
					// Ensure that the tags attribute is set in the Terraform state
					resource.TestCheckResourceAttrSet("google_dataproc_metastore_service.default", "tags.%"),
					// Perform an out-of-band check for the tag binding via API
					testAccCheckMetastoreServiceHasTagBindings(t),
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

// testAccCheckMetastoreServiceHasTagBindings verifies that a resource has tag bindings created.
// It iterates through the state, finds the metastore service, and makes an API call
// to the cloudresourcemanager tagBindings list endpoint to verify the presence of tags.
func testAccCheckMetastoreServiceHasTagBindings(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_metastore_service" {
				continue
			}
			if strings.HasPrefix(name, "resource.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			// The ID from the state file (e.g., projects/p/locations/l/services/s)
			// is the correct parent for the API call to list tag bindings.
			parentURL := fmt.Sprintf("//metastore.googleapis.com/projects/%s", rs.Primary.ID)

			// The tagBindings API endpoint is at the v3 version of the Cloud Resource Manager API.
			url := fmt.Sprintf("https://cloudresourcemanager.googleapis.com/v3/tagBindings?parent=%s", parentURL)

			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				return fmt.Errorf("Error during API request to check tag bindings: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to retrieve tag bindings. HTTP status: %s", resp.Status)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("Failed to read response body: %v", err)
			}

			var result struct {
				TagBindings []interface{} `json:"tagBindings"`
			}
			if err := json.Unmarshal(body, &result); err != nil {
				return fmt.Errorf("Failed to unmarshal JSON response: %v", err)
			}

			if len(result.TagBindings) == 0 {
				return fmt.Errorf("No tag bindings found for resource %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccMetastoreServiceTags(context map[string]interface{}) string {
	return acctest.Nprintf(`resource "google_dataproc_metastore_service" "default" {
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
}`, context)
}

func testAccCheckDataprocMetastoreServiceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dataproc_metastore_service" {
				continue
			}
			if strings.HasPrefix(name, "resource.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{MetastoreBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("MetastoreService still exists at %s", url)
			}
		}

		return nil
	}
}
