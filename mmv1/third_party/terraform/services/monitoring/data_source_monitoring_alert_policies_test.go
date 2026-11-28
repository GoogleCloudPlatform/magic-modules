package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/monitoring"
)

// testAccDataSourceGoogleMonitoringAlertPolicies_basic runs inside TestAccMonitoringAlertPolicy, since alert policy tests cannot run in parallel.
func testAccDataSourceGoogleMonitoringAlertPolicies_basic(t *testing.T) {
	alertName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	conditionName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlertPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringAlertPolicies_basicConfig(alertName, conditionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_monitoring_alert_policies.all", "id", fmt.Sprintf("projects/%s/alertPolicies", envvar.GetTestProjectFromEnv())),
					resource.TestCheckResourceAttrPair("data.google_monitoring_alert_policies.all", "project", "google_monitoring_alert_policy.basic", "project"),
					// The shared test project holds other policies, so there is no count check and the list must hold this one.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_monitoring_alert_policies.all", "alert_policies.*", map[string]string{
						"display_name":              alertName,
						"combiner":                  "OR",
						"enabled":                   "false",
						"user_labels.my-label":      "my-label-value",
						"conditions.0.display_name": conditionName,
						"conditions.0.condition_threshold.0.comparison": "COMPARISON_GT",
					}),
					resource.TestCheckTypeSetElemAttrPair("data.google_monitoring_alert_policies.all", "alert_policies.*.name", "google_monitoring_alert_policy.basic", "name"),
					resource.TestCheckTypeSetElemAttrPair("data.google_monitoring_alert_policies.all", "alert_policies.*.id", "google_monitoring_alert_policy.basic", "id"),
				),
			},
		},
	})
}

// testAccDataSourceGoogleMonitoringAlertPolicies_basicConfig creates a disabled policy and reads the project from it, so the list is read after it exists.
func testAccDataSourceGoogleMonitoringAlertPolicies_basicConfig(alertName, conditionName string) string {
	return fmt.Sprintf(`
resource "google_monitoring_alert_policy" "basic" {
  display_name = "%s"
  enabled      = false
  combiner     = "OR"

  user_labels = {
    my-label = "my-label-value"
  }

  conditions {
    display_name = "%s"

    condition_threshold {
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }

      duration        = "60s"
      comparison      = "COMPARISON_GT"
      filter          = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      threshold_value = "0.5"
    }
  }
}

data "google_monitoring_alert_policies" "all" {
  project = google_monitoring_alert_policy.basic.project
}
`, alertName, conditionName)
}
