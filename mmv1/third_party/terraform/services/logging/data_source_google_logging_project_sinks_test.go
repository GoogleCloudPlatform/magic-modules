package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/logging"
	_ "github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func TestAccDataSourceGoogleLoggingProjectSinks_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name": envvar.GetTestProjectFromEnv(),
		"sink_name":    "tf-test-sinks-ds-" + acctest.RandString(t, 10),
		"bucket_name":  "tf-test-sinks-ds-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleLoggingProjectSinks_basic(context),
				Check: resource.ComposeTestCheckFunc(
					// The shared test project holds other sinks, so there is no count check.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_logging_project_sinks.all", "sinks.*", map[string]string{
						"name":        context["sink_name"].(string),
						"destination": "storage.googleapis.com/" + context["bucket_name"].(string),
						"filter":      fmt.Sprintf(`logName="projects/%s/logs/compute.googleapis.com%%2Factivity_log" AND severity>=ERROR`, context["project_name"]),
						"disabled":    "false",
					}),
					resource.TestCheckTypeSetElemAttrPair("data.google_logging_project_sinks.all", "sinks.*.writer_identity", "google_logging_project_sink.basic", "writer_identity"),
					resource.TestCheckTypeSetElemNestedAttrs("data.google_logging_project_sinks.all", "sinks.*", map[string]string{"name": "_Required"}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleLoggingProjectSinks_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_logging_project_sink" "basic" {
  name        = "%{sink_name}"
  project     = "%{project_name}"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%{project_name}/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%{bucket_name}"
  location = "US"
}

data "google_logging_project_sinks" "all" {
  project = google_logging_project_sink.basic.project
}
`, context)
}
