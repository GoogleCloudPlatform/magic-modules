package monitoring_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccMonitoringSnooze_monitoringSnoozeBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGoogleMonitoringSnoozeWasCancelledAndRemovedFromState(t, "google_monitoring_snooze.snooze"),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSnooze_monitoringSnoozeBasicExample(context),
			},
			{
				ResourceName:      "google_monitoring_snooze.snooze",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testGoogleMonitoringSnooze_removed(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleMonitoringSnoozeWasCancelledAndRemovedFromState(t, "google_monitoring_snooze.snooze"),
				),
			},
		},
	})
}

func testAccMonitoringSnooze_monitoringSnoozeBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_alert_policy" "tf_test_alert_policy%{random_suffix}" {
  display_name = "My Alert Policy%{random_suffix}"
  combiner     = "OR"
  conditions {
    display_name = "test condition"
    condition_threshold {
      filter     = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }
}

resource "google_monitoring_snooze"  "snooze" {
  display_name = "My Snooze%{random_suffix}"
  
  interval {
    start_time = replace(timeadd(plantimestamp(), "24h"), "/T.*/", "T00:00:00Z")
    end_time   = replace(timeadd(plantimestamp(), "25h"), "/T.*/", "T01:00:00Z")
  }

  criteria {
    policies = [
        google_monitoring_alert_policy.tf_test_alert_policy%{random_suffix}.id
    ]
  }
}
`, context)
}

// Monitoring Snooze cannot be deleted. This ensures that the Monitoring Snooze resource was removed from state,
// even though the server-side resource was not removed.
func testAccCheckGoogleMonitoringSnoozeWasCancelledAndRemovedFromState(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_snooze" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			_, ok := s.RootModule().Resources[resourceName]
			if ok {
				return fmt.Errorf("Resource was not removed from state: %s", resourceName)
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{MonitoringBasePath}}v3/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err != nil {
				return errwrap.Wrapf("Failed to get expected Monitoring Snooze resource", err)
			}

			// It would be great to test whether the snooze has been cancelled here however,
			// that state of of a Snooze is not currently exposed through the Cloud Monitoring API
			// if response["state"] != "Cancelled" {
			// 	return fmt.Errorf("Monitoring Snooze has not been cancelled")
			// }
		}

		return nil
	}
}

func testGoogleMonitoringSnooze_removed(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_alert_policy" "tf_test_alert_policy%{random_suffix}" {
  display_name = "My Alert Policy%{random_suffix}"
  combiner     = "OR"
  conditions {
    display_name = "test condition"
    condition_threshold {
      filter     = "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\""
      duration   = "60s"
      comparison = "COMPARISON_GT"
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_RATE"
      }
    }
  }
}
`, context)
}
