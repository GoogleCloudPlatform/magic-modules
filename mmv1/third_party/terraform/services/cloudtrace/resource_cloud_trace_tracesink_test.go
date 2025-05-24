package cloudtrace_test

import (
	"fmt"
)

func testAccCloudTraceTraceSink_basic(name string) string {
	return fmt.Sprintf(`
resource "google_cloud_trace_trace_sink" "default" {
  name = "%s"
  output_config {
    destination = "bigquery.googleapis.com/projects/my-project/datasets/my-dataset"
  }
}
`, name)
}

func testAccCloudTraceTraceSink_update(name string) string {
	return fmt.Sprintf(`
resource "google_cloud_trace_trace_sink" "default" {
  name = "%s"
  output_config {
    destination = "bigquery.googleapis.com/projects/my-project/datasets/my-dataset-2"
  }
`, name)
}
