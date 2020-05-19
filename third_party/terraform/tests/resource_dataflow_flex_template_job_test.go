package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataflowFlexTemplateJob_simple(t *testing.T) {
	t.Parallel()

	randStr := randString(t, 10)
	bucket := "tf-test-dataflow-gcs-" + randStr
	job := "tf-test-dataflow-job-" + randStr

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataflowJobDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataflowFlowFlexTemplateJob_basic(bucket, job),
				Check: resource.ComposeTestCheckFunc(
					testAccDataflowJobExists(t, "google_dataflow_job.big_data"),
				),
			},
		},
	})
}

func testAccDataflowFlowFlexTemplateJob_basic(bucket, job string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "temp" {
  name = "%s"
  force_destroy = true
}

resource "google_storage_bucket_object" "flex_template" {
	name   = "flex_template.json"
	bucket = "%s"
	content = "%s"
}

resource "google_dataflow_flex_template_job" "big_data" {
  name = "%s"
  container_spec_gcs_path = "%s"
  on_delete = "cancel"
}
`, bucket, bucket, flexTemplateContent(), bucket, job)
}

func flexTemplateContent() string {
	return `
{
	"name": "Streaming Beam SQL",
	"description": "An Apache Beam streaming pipeline that reads JSON encoded messages from Pub/Sub, uses Beam SQL to transform the message data, and writes the results to a BigQuery",
	"parameters": [
		{
		"name": "inputSubscription",
		"label": "Pub/Sub input subscription.",
		"helpText": "Pub/Sub subscription to read from.",
		"regexes": [
			"[-_.a-zA-Z0-9]+"
		]
		},
		{
		"name": "outputTable",
		"label": "BigQuery output table",
		"helpText": "BigQuery table spec to write to, in the form 'project:dataset.table'.",
		"is_optional": true,
		"regexes": [
			"[^:]+:[^.]+[.].+"
		]
		}
	]
}
`
}
