package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMonitoringMetricDescriptor_update(t *testing.T) {

	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMetricDescriptorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMetricDescriptor_update("initial description", "initial display name", "30s", "30s"),
			},
			{
				ResourceName:            "google_monitoring_metric_descriptor.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata", "launch_stage"},
			},
			{
				Config: testAccMonitoringMetricDescriptor_update("updated description", "updated display name", "60s", "60s"),
			},
			{
				ResourceName:            "google_monitoring_metric_descriptor.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata", "launch_stage"},
			},
			{
				Config: testAccMonitoringMetricDescriptor_omittedFields(),
			},
			{
				ResourceName:            "google_monitoring_metric_descriptor.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata", "launch_stage", "description", "display_name"},
			},
		},
	})
}

func testAccMonitoringMetricDescriptor_update(description, displayName, samplePeriod, ingestDelay string) string {
	return fmt.Sprintf(`
resource "google_monitoring_metric_descriptor" "basic" {
	description = "%s"
	display_name = "%s"
	type = "custom.googleapis.com/stores/daily_sales"
	metric_kind = "GAUGE"
	value_type = "DOUBLE"
	unit = "{USD}"
	labels {
		key = "key"
		value_type = "STRING"
		description = "description"
	}
	launch_stage = "BETA"
	metadata {
		sample_period = "%s"
		ingest_delay = "%s"
	}
}
`, description, displayName, samplePeriod, ingestDelay,
	)
}

func testAccMonitoringMetricDescriptor_omittedFields() string {
	return `
resource "google_monitoring_metric_descriptor" "basic" {
	type = "custom.googleapis.com/stores/daily_sales"
	metric_kind = "GAUGE"
	value_type = "DOUBLE"
	unit = "{USD}"
	labels {
		key = "key"
		value_type = "STRING"
		description = "description"
	}
	launch_stage = "BETA"
	metadata {
		sample_period = "30s"
		ingest_delay = "30s"
	}
}
`
}
