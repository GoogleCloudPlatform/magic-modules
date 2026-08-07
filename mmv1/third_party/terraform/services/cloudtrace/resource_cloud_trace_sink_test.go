package cloudtrace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCloudTraceSink_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudTraceSink_basic(context),
			},
			{
				ResourceName:      "google_cloud_trace_sink.sink",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update the export destination to a different dataset.
				Config: testAccCloudTraceSink_updated(context),
			},
			{
				ResourceName:      "google_cloud_trace_sink.sink",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudTraceSink_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_bigquery_dataset" "first" {
  dataset_id = "tf_test_trace_sink_first_%{random_suffix}"
  location   = "US"
}

resource "google_bigquery_dataset" "second" {
  dataset_id = "tf_test_trace_sink_second_%{random_suffix}"
  location   = "US"
}

resource "google_cloud_trace_sink" "sink" {
  project = data.google_project.project.number
  sink_id = "tf-test-trace-sink-%{random_suffix}"

  output_config {
    destination = "bigquery.googleapis.com/projects/${data.google_project.project.number}/datasets/${google_bigquery_dataset.first.dataset_id}"
  }
}
`, context)
}

func testAccCloudTraceSink_updated(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_bigquery_dataset" "first" {
  dataset_id = "tf_test_trace_sink_first_%{random_suffix}"
  location   = "US"
}

resource "google_bigquery_dataset" "second" {
  dataset_id = "tf_test_trace_sink_second_%{random_suffix}"
  location   = "US"
}

resource "google_cloud_trace_sink" "sink" {
  project = data.google_project.project.number
  sink_id = "tf-test-trace-sink-%{random_suffix}"

  output_config {
    destination = "bigquery.googleapis.com/projects/${data.google_project.project.number}/datasets/${google_bigquery_dataset.second.dataset_id}"
  }
}
`, context)
}
