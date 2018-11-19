package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMonitoringUptimeCheckConfig_update(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringUptimeCheckConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringUptimeCheckConfig_update(acctest.RandString(10)),
			},
			{
				ResourceName:      "google_monitoring_uptime_check_config.http",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func partialHttpCheck(port, user, pass string) string {
	auth := ""
	if user != "" || pass != "" {
		auth = fmt.Sprintf(`
    auth_ino = {
      username = "%s"
      password = "%s"
		}
		`, user, pass)
	}

	return fmt.Sprintf(`
resource "google_monitoring_uptime_check_config" "http" {
  display_name = "http-uptime-check-%s"
  timeout = "60s"

  http_check = {
    path = "/some-path"
		port = "%s"
		%s
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = "chrisst-TODO"
      host = "192.168.1.1"
    }
  }

  content_matchers = {
    content = "example"
  }
}
`, acctest.RandString(4), port, auth,
	)
}

func testAccMonitoringUptimeCheckConfig_update(val string) string {
	return fmt.Sprintf(`
resource "google_monitoring_uptime_check_config" "http" {
  display_name = "http-uptime-check-%s"
  timeout = "60s"

  http_check = {
    path = "/some-path"
    port = "8010"
    auth_ino = {
      username = "name"
      password = "basic_auth"
    }
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = "chrisst-TODO"
      host = "192.168.1.1"
    }
  }

  content_matchers = {
    content = "example"
  }
}
`, val,
	)
}

// Error 400: 'metric' is not a valid resource type. Only 'gce_instance', 'aws_ec2_instance', 'aws_elb_load_balancer', 'gae_app', and 'uptime_url' are valid.
