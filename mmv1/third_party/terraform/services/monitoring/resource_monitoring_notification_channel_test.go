package monitoring_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMonitoringNotificationChannel_update(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_update("email", `email_address = "fake_email@blahblah.com"`, "true"),
			},
			{
				ResourceName:      "google_monitoring_notification_channel.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringNotificationChannel_update("sms", `number = "+16502530000"`, "false"),
			},
			{
				ResourceName:      "google_monitoring_notification_channel.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMonitoringNotificationChannel_updateLabels_slack(t *testing.T) {
	// Slack auth_token required for test not to fail, skipping test till interal testing slack can be created
	t.Skip()
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_updateLabels_slack(),
			},
			{
				ResourceName:            "google_monitoring_notification_channel.slack",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.auth_token"},
			},
		},
	})
}

func TestAccMonitoringNotificationChannel_updateLabels(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_updateLabels(),
			},
			{
				ResourceName:            "google_monitoring_notification_channel.pagerduty",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.service_key", "sensitive_labels"},
			},
			{
				ResourceName:            "google_monitoring_notification_channel.basicauth",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.password", "sensitive_labels"},
			},
		},
	})
}

func TestAccMonitoringNotificationChannel_updateSensitiveLabels_slack(t *testing.T) {
	// Slack auth_token required for test not to fail, skipping test till interal testing slack can be created
	t.Skip()
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_updateSensitiveLabels_slack(),
			},
			// sensitive labels for notification channels are either obfuscated or not returned by the upstream
			// API. Therefore when re-importing a resource we cannot know what the value is.
			{
				ResourceName:            "google_monitoring_notification_channel.slack",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "sensitive_labels.auth_token"},
			},
		},
	})
}
func TestAccMonitoringNotificationChannel_updateSensitiveLabels(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_updateSensitiveLabels(),
			},
			// sensitive labels for notification channels are either obfuscated or not returned by the upstream
			// API. Therefore when re-importing a resource we cannot know what the value is.
			{
				ResourceName:            "google_monitoring_notification_channel.pagerduty",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.service_key", "sensitive_labels"},
			},
			{
				ResourceName:            "google_monitoring_notification_channel.basicauth",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.password", "sensitive_labels"},
			},
		},
	})
}

func TestAccMonitoringNotificationChannel_updateSensitiveLabelsWriteOnly_slack(t *testing.T) {
	// Slack auth_token required for test not to fail, skipping test till internal testing slack can be created
	t.Skip()
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_sensitiveLabelsWriteOnly_slack("token1", 1),
			},
			{
				ResourceName:            "google_monitoring_notification_channel.slack",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.auth_token", "sensitive_labels"},
			},
			{
				Config: testAccMonitoringNotificationChannel_sensitiveLabelsWriteOnly_slack("token2", 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_monitoring_notification_channel.slack", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_monitoring_notification_channel.slack",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.auth_token", "sensitive_labels"},
			},
		},
	})
}

func TestAccMonitoringNotificationChannel_updateSensitiveLabelsWriteOnly(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringNotificationChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringNotificationChannel_sensitiveLabelsWriteOnly("key1", "pass1", 1),
			},
			{
				ResourceName:            "google_monitoring_notification_channel.pagerduty",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.service_key", "sensitive_labels"},
			},
			{
				ResourceName:            "google_monitoring_notification_channel.basicauth",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.password", "sensitive_labels"},
			},
			{
				Config: testAccMonitoringNotificationChannel_sensitiveLabelsWriteOnly("key2", "pass2", 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_monitoring_notification_channel.pagerduty", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("google_monitoring_notification_channel.basicauth", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_monitoring_notification_channel.pagerduty",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.service_key", "sensitive_labels"},
			},
			{
				ResourceName:            "google_monitoring_notification_channel.basicauth",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels.%", "labels.password", "sensitive_labels"},
			},
		},
	})
}

func testAccMonitoringNotificationChannel_update(channel, labels, enabled string) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "update" {
  display_name = "IntTest Notification Channel"
  type         = "%s"
  labels = {
    %s
  }

  enabled = "%s"
}
`, channel, labels, enabled,
	)
}

func testAccMonitoringNotificationChannel_updateLabels_slack() string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "slack" {
	display_name = "TFTest Slack Channel"
	type         = "slack"
	labels = {
		"auth_token"   = "one"
		"channel_name" = "#foobar"
	}
}
`)
}

func testAccMonitoringNotificationChannel_updateLabels() string {
	return fmt.Sprintf(`


resource "google_monitoring_notification_channel" "basicauth" {
	display_name = "TFTest Basicauth Channel"
	type         = "webhook_basicauth"
	labels = {
		"password" = "somepassword"
		"username" = "username"
		"url"      = "http://fakeurl.com"
	}
}

resource "google_monitoring_notification_channel" "pagerduty" {
	display_name = "TFTest Pagerduty Channel"
	type         = "pagerduty"
	labels = {
		"service_key" = "some_service_key"
	}
}
`)
}

func testAccMonitoringNotificationChannel_updateSensitiveLabels_slack() string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "slack" {
	display_name = "TFTest Slack Channel"
	type         = "slack"
	labels = {
		"channel_name" = "#foobar"
	}

	sensitive_labels {
		auth_token = "one"
	}
}
`)
}

func testAccMonitoringNotificationChannel_updateSensitiveLabels() string {
	return fmt.Sprintf(`

resource "google_monitoring_notification_channel" "basicauth" {
	display_name = "TFTest Basicauth Channel"
	type         = "webhook_basicauth"
	labels = {
		"username" = "username"
		"url"      = "http://fakeurl.com"
	}

	sensitive_labels {
		password = "somepassword"
	}
}

resource "google_monitoring_notification_channel" "pagerduty" {
	display_name = "TFTest Pagerduty Channel"
	type         = "pagerduty"

	sensitive_labels {
		service_key = "some_service_key"
	}
}
`)
}

func testAccMonitoringNotificationChannel_sensitiveLabelsWriteOnly_slack(authToken string, version int) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "slack" {
	display_name = "TFTest Slack Channel"
	type         = "slack"
	labels = {
		"channel_name" = "#foobar"
	}

	sensitive_labels {
		auth_token_wo         = "%s"
		auth_token_wo_version = %d
	}
}
`, authToken, version)
}

func testAccMonitoringNotificationChannel_sensitiveLabelsWriteOnly(serviceKey, password string, version int) string {
	return fmt.Sprintf(`
resource "google_monitoring_notification_channel" "basicauth" {
	display_name = "TFTest Basicauth Channel"
	type         = "webhook_basicauth"
	labels = {
		"username" = "username"
		"url"      = "http://fakeurl.com"
	}

	sensitive_labels {
		password_wo         = "%s"
		password_wo_version = %d
	}
}

resource "google_monitoring_notification_channel" "pagerduty" {
	display_name = "TFTest Pagerduty Channel"
	type         = "pagerduty"

	sensitive_labels {
		service_key_wo         = "%s"
		service_key_wo_version = %d
	}
}
`, password, version, serviceKey, version)
}
