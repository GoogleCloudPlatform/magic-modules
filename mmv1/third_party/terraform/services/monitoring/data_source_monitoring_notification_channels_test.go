package monitoring_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/monitoring"
)

func TestAccDataSourceGoogleMonitoringNotificationChannels_basic(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	password := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 16))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleMonitoringNotificationChannels_basicConfig(displayName, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_monitoring_notification_channels.all", "id", fmt.Sprintf("projects/%s/notificationChannels", envvar.GetTestProjectFromEnv())),
					resource.TestCheckResourceAttrPair("data.google_monitoring_notification_channels.all", "project", "google_monitoring_notification_channel.basic", "project"),
					// The shared test project holds other channels, so there is no count check and the list must hold these two.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_monitoring_notification_channels.all", "notification_channels.*", map[string]string{
						"display_name":         displayName,
						"description":          "A channel for the notification channels data source test",
						"type":                 "email",
						"enabled":              "false",
						"labels.email_address": fmt.Sprintf("%s@google.com", displayName),
						"user_labels.testname": displayName,
					}),
					resource.TestCheckTypeSetElemAttrPair("data.google_monitoring_notification_channels.all", "notification_channels.*.name", "google_monitoring_notification_channel.basic", "name"),
					resource.TestCheckTypeSetElemAttrPair("data.google_monitoring_notification_channels.all", "notification_channels.*.id", "google_monitoring_notification_channel.basic", "id"),
					resource.TestCheckTypeSetElemNestedAttrs("data.google_monitoring_notification_channels.all", "notification_channels.*", map[string]string{
						"display_name":    displayName + "-basicauth",
						"type":            "webhook_basicauth",
						"labels.username": "username",
						"labels.url":      "http://fakeurl.com",
					}),
					testAccCheckDataSourceGoogleMonitoringNotificationChannelsHideSecret("data.google_monitoring_notification_channels.all", password),
				),
			},
		},
	})
}

// testAccCheckDataSourceGoogleMonitoringNotificationChannelsHideSecret fails when an attribute of the data source holds the secret in clear.
func testAccCheckDataSourceGoogleMonitoringNotificationChannelsHideSecret(dataSourceName, secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", dataSourceName)
		}
		for key, value := range ds.Primary.Attributes {
			if strings.Contains(value, secret) {
				return fmt.Errorf("%s holds the secret in clear, want it obfuscated or left out by the API", key)
			}
		}
		return nil
	}
}

// testAccDataSourceGoogleMonitoringNotificationChannels_basicConfig creates an email channel and a channel with a password, and lists the project once both exist.
func testAccDataSourceGoogleMonitoringNotificationChannels_basicConfig(displayName, password string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "basic" {
  display_name = "%s"
  description  = "A channel for the notification channels data source test"
  type         = "email"
  enabled      = false

  labels = {
    email_address = "%s@google.com"
  }

  user_labels = {
    testname = "%s"
  }
}

resource "google_monitoring_notification_channel" "basicauth" {
  display_name = "%s-basicauth"
  type         = "webhook_basicauth"

  labels = {
    username = "username"
    url      = "http://fakeurl.com"
  }

  sensitive_labels {
    password = "%s"
  }
}

data "google_monitoring_notification_channels" "all" {
  project    = google_monitoring_notification_channel.basic.project
  depends_on = [google_monitoring_notification_channel.basicauth]
}
`, displayName, displayName, displayName, displayName, password)
}
