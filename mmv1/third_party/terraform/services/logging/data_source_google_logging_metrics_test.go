package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/logging"
)

func TestAccDataSourceGoogleLoggingMetrics_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-metrics-ds-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingMetricDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleLoggingMetricsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					// The shared test project holds other metrics, so there is no count check.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_logging_metrics.all", "metrics.*", map[string]string{
						"name":                            name,
						"filter":                          "resource.type=gae_app AND severity>=ERROR",
						"disabled":                        "false",
						"value_extractor":                 "EXTRACT(jsonPayload.request)",
						"label_extractors.mass":           "EXTRACT(jsonPayload.request)",
						"metric_descriptor.0.metric_kind": "DELTA",
						"metric_descriptor.0.value_type":  "DISTRIBUTION",
						"metric_descriptor.0.labels.#":    "2",
						"bucket_options.0.linear_buckets.0.num_finite_buckets": "3",
					}),
				),
			},
		},
	})
}

// testAccDataSourceGoogleLoggingMetricsConfig creates a distribution metric, which fills every kind of nested field.
func testAccDataSourceGoogleLoggingMetricsConfig(name string) string {
	return fmt.Sprintf(`
resource "google_logging_metric" "foobar" {
  name   = "%s"
  filter = "resource.type=gae_app AND severity>=ERROR"

  metric_descriptor {
    metric_kind = "DELTA"
    value_type  = "DISTRIBUTION"
    unit        = "1"
    labels {
      key         = "mass"
      value_type  = "STRING"
      description = "amount of matter"
    }
    labels {
      key         = "sku"
      value_type  = "INT64"
      description = "Identifying number for item"
    }
  }

  value_extractor = "EXTRACT(jsonPayload.request)"
  label_extractors = {
    "mass" = "EXTRACT(jsonPayload.request)"
    "sku"  = "EXTRACT(jsonPayload.id)"
  }

  bucket_options {
    linear_buckets {
      num_finite_buckets = 3
      width              = 1
      offset             = 1
    }
  }
}

data "google_logging_metrics" "all" {
  project = google_logging_metric.foobar.project
}
`, name)
}
