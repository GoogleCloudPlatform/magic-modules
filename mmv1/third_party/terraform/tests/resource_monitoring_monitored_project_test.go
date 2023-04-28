package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMonitoringMonitoredProject_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        GetTestOrgFromEnv(t),
		"project_id":    GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMonitoredProject_basic(context),
			},
			{
				ResourceName:            "google_monitoring_monitored_project.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metrics_scope"},
			},
		},
	})
}

func testAccMonitoringMonitoredProject_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "%{project_id}"
  name          = google_project.basic.name
}

resource "google_project" "basic" {
  project_id = "tf-test-m-id%{random_suffix}"
  name       = "tf-test-m-id%{random_suffix}"
  org_id     = "%{org_id}"
}
`, context)
}

func testAccCheckMonitoringMonitoredProjectDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_monitored_project" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("MonitoringMonitoredProject still exists at %s", url)
			}
		}

		return nil
	}
}
