package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccMonitoringMetricDescriptor_initialConfig("initial description", "initial display name"),
			},
			{
				ResourceName:            "google_monitoring_metric_descriptor.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata", "launch_stage"},
			},
			{
				Config: testAccMonitoringMetricDescriptor_updatedConfig("updated description", "updated display name"),
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

func testAccMonitoringMetricDescriptor_initialConfig(description, displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_metric_descriptor" "basic" {
	description = "%s"
	display_name = "%s"
	type = "custom.googleapis.com/stores/daily_sales"
	metric_kind = "GAUGE"
	value_type = "DOUBLE"
	unit = "{USD}"
	labels {
		key = "key1"
		value_type = "STRING"
		description = "description1"
	}
	launch_stage = "BETA"
	metadata {
		sample_period = "30s"
		ingest_delay = "30s"
	}
}
`, description, displayName)
}

func testAccMonitoringMetricDescriptor_updatedConfig(description, displayName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_metric_descriptor" "basic" {
	description = "%s"
	display_name = "%s"
	type = "custom.googleapis.com/stores/daily_sales"
	metric_kind = "GAUGE"
	value_type = "DOUBLE"
	unit = "{USD}"
	labels {
		key = "key1"
		value_type = "STRING"
		description = "description1"
	}
	launch_stage = "BETA"
	metadata {
		sample_period = "30s"
		ingest_delay = "30s"
	}
}
`, description, displayName)
}

func testAccMonitoringMetricDescriptor_omittedFields() string {
	return `
resource "google_monitoring_metric_descriptor" "basic" {
	type = "custom.googleapis.com/stores/daily_sales"
	metric_kind = "GAUGE"
	value_type = "DOUBLE"
	unit = "{USD}"
	labels {
		key = "key1"
		value_type = "STRING"
		description = "description1"
	}
	launch_stage = "BETA"
	metadata {
		sample_period = "30s"
		ingest_delay = "30s"
	}
}
`
}
