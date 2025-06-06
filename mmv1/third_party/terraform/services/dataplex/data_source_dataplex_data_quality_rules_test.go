package dataplex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccDataplexDataQualityRules(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
		"location":      envvar.GetTestRegionFromEnv(),
		"data_scan_id":  "tf-test-datascan-profile-id",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDataQualityRules_datascan_config(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dataplex_datascan.tf_test_datascan_profile", "project", context["project"].(string)),
					resource.TestCheckResourceAttr("google_dataplex_datascan.tf_test_datascan_profile", "location", context["location"].(string)),
					resource.TestCheckResourceAttr("google_dataplex_datascan.tf_test_datascan_profile", "data_scan_id", context["data_scan_id"].(string)),
					resource.TestCheckResourceAttr("google_dataplex_datascan.tf_test_datascan_profile", "data.0.resource", "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/samples/tables/shakespeare"),
				),
			},
			{
				ResourceName:            "google_dataplex_datascan.tf_test_datascan_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
			{
				RefreshState: true,
				Check:        testAccDataplexDataScanJobTriggerRunAndWaitUntilComplete(t, "google_dataplex_datascan.tf_test_datascan_profile"),
			},
			{
				ResourceName:            "google_dataplex_datascan.tf_test_datascan_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "execution_status"},
			},
			{
				Config: testAccDataplexDataQualityRules_rules_config(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dataplex_data_quality_rules.generated_dq_rules", "rules.#", "7"),
				),
			},
		},
	})
}

func testAccDataplexDataQualityRules_datascan_config(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_dataplex_datascan" "tf_test_datascan_profile" {
			location     = "%{location}"
			data_scan_id = "%{data_scan_id}"

			data {
				resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/samples/tables/shakespeare"
			}

			execution_spec {
				trigger {
					on_demand {}
				}
			}

			data_profile_spec {}

			project = "%{project}"
		}`, context)
}

func testAccDataplexDataScanJobTriggerRunAndWaitUntilComplete(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		t.Logf("DEBUG: Attributes for resource: %+v", rs.Primary.Attributes)
		// t.Logf("DEBUG: resourceName: %s", resourceName)
		// t.Logf("DEBUG: rr: %s", "google_dataplex_datascan.tf_test_datascan_profile.data_scan_id")

		config := acctest.GoogleProviderConfig(t)
		url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DataplexBasePath}}projects/{{project}}/locations/{{location}}/dataScans/{{data_scan_id}}:run")
		if err != nil {
			return fmt.Errorf("Failed to generate URL for triggering datascan run: %s", err)
		}

		billingProject := ""

		if config.BillingProject != "" {
			billingProject = config.BillingProject
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: config.UserAgent,
		})

		if err != nil {
			return fmt.Errorf("Request for triggering data scan run failed: %s", err)
		}

		t.Logf("res[\"job\"]: %s", res["job"])

		dataScanJobId := extractDataScanJobId(res["job"])
		t.Logf("data scan job id: %s", dataScanJobId)
		dataScanJobState := extractDataScanJobState(res["job"])
		t.Logf("data scan job state: %s", dataScanJobState)

		for dataScanJobState != "SUCCEEDED" {
			dataScanJobState, err = getDataScanJobState(t, rs, dataScanJobId)
			if err != nil {
				return fmt.Errorf("Getting data scan job state failed: failed to get state: %s", err)
			}

			switch dataScanJobState {
			case "STATE_UNSPECIFIED", "RUNNING", "PENDING":
				t.Logf("Data scan job stateeee: %s, waiting for the job to finish", dataScanJobState)
				time.Sleep(10 * time.Second) // Pause for 10 seconds
			case "CANCELING", "CANCELLED", "FAILED":
				return fmt.Errorf("Data scan job failed: Invalid state: %s", dataScanJobState)
			case "SUCCEEDED":
				t.Logf("Data scan job state SUCCEEDED: %s", dataScanJobState)
				return nil
			default:
				return fmt.Errorf("Getting data scan job state failed: invalid state: %s", dataScanJobState)
			}
		}
		time.Sleep(10 * time.Second) // Pause for 10 seconds

		return nil
	}
}

func testAccDataplexDataQualityRules_rules_config(context map[string]interface{}) string {
	return acctest.Nprintf(`
		resource "google_dataplex_datascan" "tf_test_datascan_profile" {
			location     = "%{location}"
			data_scan_id = "%{data_scan_id}"

			data {
				resource = "//bigquery.googleapis.com/projects/bigquery-public-data/datasets/samples/tables/shakespeare"
			}

			execution_spec {
				trigger {
					on_demand {}
				}
			}

			data_profile_spec {}

			project = "%{project}"
		}
			
		data "google_dataplex_data_quality_rules" "generated_dq_rules" {
			project		 = google_dataplex_datascan.tf_test_datascan_profile.project
			location	 = google_dataplex_datascan.tf_test_datascan_profile.location
			data_scan_id = google_dataplex_datascan.tf_test_datascan_profile.data_scan_id
		}`, context)
}

/*
 Error: Error when reading or editing DataQualityRules "projects/google_dataplex_datascan.tf_test_datascan_profile.project/locations/google_dataplex_datascan.tf_test_datascan_profile.location/dataScans/google_dataplex_datascan.tf_test_datascan_profile.data_scan_id": googleapi: Error 403: Permission denied on resource project google_dataplex_datascan.tf_test_datascan_profile.project.
        Details:
        [
          {
            "@type": "type.googleapis.com/google.rpc.ErrorInfo",
            "domain": "googleapis.com",
            "metadata": {
              "consumer": "projects/google_dataplex_datascan.tf_test_datascan_profile.project",
              "containerInfo": "google_dataplex_datascan.tf_test_datascan_profile.project",
              "service": "dataplex.googleapis.com"
            },
            "reason": "CONSUMER_INVALID"
          },

*/

/*
func testAccDataplexDataQualityRules_rules_config(context map[string]interface{}) string {
	return acctest.Nprintf(`
		data "google_dataplex_data_quality_rules" "generated_dq_rules" {
			project		 = "google_dataplex_datascan.tf_test_datascan_profile.project"
			location	 = "google_dataplex_datascan.tf_test_datascan_profile.location"
			data_scan_id = "google_dataplex_datascan.tf_test_datascan_profile.data_scan_id"
		}`, context)
}
*/

/*
googleapi: Error 400: Provided DataScan 'projects/1044284642890/locations/us-central1/dataScans/tf-test-datascan-profile-id' does not exist.

	Details:
	[
	  {
	    "@type": "type.googleapis.com/google.rpc.DebugInfo",
	    "detail": "generic::invalid_argument: Provided DataScan 'projects/1044284642890/locations/us-central1/dataScans/tf-test-datascan-profile-id' does not exist.",
*/

// func testAccDataplexDataQualityRules_rules_config(context map[string]interface{}) string {
// 	return acctest.Nprintf(`
// 		data "google_dataplex_data_quality_rules" "generated_dq_rules" {
// 			project		 = "hvaidya-test-dataplex"
// 			location	 = "%{location}"
// 			data_scan_id = "%{data_scan_id}"
// 		}`, context)
// }

func getDataScanJobState(t *testing.T, rs *terraform.ResourceState, dataScanJobId string) (string, error) {
	config := acctest.GoogleProviderConfig(t)
	// GET https://dataplex.googleapis.com/v1/{name=projects/*/locations/*/dataScans/*/jobs/*}
	url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{DataplexBasePath}}projects/{{project}}/locations/{{location}}/dataScans/{{data_scan_id}}/jobs/"+dataScanJobId)
	if err != nil {
		return "", fmt.Errorf("Failed to generate URL for getting data scan job state: %s", err)
	}
	t.Logf("DEBUG: url: %s", url)

	billingProject := ""

	if config.BillingProject != "" {
		billingProject = config.BillingProject
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: config.UserAgent,
	})

	if err != nil {
		return "", fmt.Errorf("Request for getting data scan job state failed: %s", err)
	}

	return extractDataScanJobState(res), nil
}

func extractDataScanJobState(job interface{}) string {
	dataScanJob := job.(map[string]interface{})
	return dataScanJob["state"].(string)
}

func extractDataScanJobId(job interface{}) string {
	dataScanJob := job.(map[string]interface{})
	return dataScanJob["uid"].(string)
}
