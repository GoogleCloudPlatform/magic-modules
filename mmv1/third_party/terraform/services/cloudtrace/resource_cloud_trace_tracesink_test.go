package cloudtrace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func testAccCloudTraceTraceSink_basic(name string, datasetID string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_cloud_trace_trace_sink" "default" {
  name = "%s"
  output_config {
    destination = "bigquery.googleapis.com/projects/my-project/datasets/%s"
  }
}
`, name, datasetID)
}

func testAccCloudTraceTraceSink_update(name string, datasetID string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_cloud_trace_trace_sink" "default" {
  name = "%s"
  output_config {
    destination = "bigquery.googleapis.com/projects/my-project/datasets/%s"
  }
}
`, name, datasetID)
}

func testAccBigQueryDataset(datasetID string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "traces"
  description                     = "Test dataset for Cloud Trace Sink"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
}
`, datasetID)
}

func TestAccCloudTraceSink_update(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset("traces"),
			},
			{
				Config: testAccCloudTraceTraceSink_basic("example", "traces"),
			},
			{
				Config: testAccCloudTraceTraceSink_update("updated example", "traces"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_cloud_trace_sink.name", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}
