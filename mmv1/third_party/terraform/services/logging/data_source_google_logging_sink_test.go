package logging_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/logging"
)

func TestAccDataSourceGoogleLoggingSink_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name": envvar.GetTestProjectFromEnv(),
		"sink_name":    "tf-test-sink-ds-" + acctest.RandString(t, 10),
		"bucket_name":  "tf-test-sink-ds-bucket-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleLoggingSink_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_logging_sink.basic",
						"google_logging_project_sink.basic",
						map[string]struct{}{
							"project":                {},
							"unique_writer_identity": {},
						},
					),
				),
			},
		},
	})
}

func TestUnitLoggingSink_OptionalSurroundingSpacesSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"surrounding spaces": {
			Old:                "value",
			New:                " value ",
			ExpectDiffSuppress: true,
		},
		"no surrounding spaces": {
			Old:                "value",
			New:                "value",
			ExpectDiffSuppress: true,
		},
		"one space each": {
			Old:                " value",
			New:                "value ",
			ExpectDiffSuppress: true,
		},
		"different values": {
			Old:                " different",
			New:                "values ",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if logging.OptionalSurroundingSpacesSuppress("filter", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func testAccDataSourceGoogleLoggingSink_basic(context map[string]interface{}) string {
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

data "google_logging_sink" "basic" {
  id = google_logging_project_sink.basic.id
}
`, context)
}
