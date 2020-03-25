package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCloudIoTRegistry_update(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistryBasic(registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTRegistryExtended(registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTRegistryBasic(registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudIoTRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudiot_registry" {
			continue
		}
		config := testAccProvider.Meta().(*Config)
		registry, _ := config.clientCloudIoT.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if registry != nil {
			return fmt.Errorf("Registry still present")
		}
	}
	return nil
}

func testAccCloudIoTRegistryBasic(registryName string) string {
	return fmt.Sprintf(`
resource "google_cloud_iot_device_registry" "%s" {
  name = "%s"
}
`, registryName, registryName)
}

func testAccCloudIoTRegistryExtended(registryName string) string {
	return fmt.Sprintf(`

resource "google_project_service" "cloud-iot-apis" {
  service = "cloudiot.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "pubsub-apis" {
  service = "pubsub.googleapis.com"

  disable_dependent_services = true
}

resource "google_pubsub_topic" "default-devicestatus" {
  name = "psregistry-test-devicestatus-%s"

  depends_on = [
    google_project_service.pubsub-apis
  ]
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "psregistry-test-telemetry-%s"

  depends_on = [
    google_project_service.pubsub-apis
  ]
}

resource "google_pubsub_topic" "additional-telemetry" {
  name = "psregistry-additional-test-telemetry-%s"

  depends_on = [
    google_project_service.pubsub-apis
  ]
}

resource "google_cloud_iot_device_registry" "%s" {
  name     = "%s"

  depends_on = [
    google_pubsub_topic.default-telemetry,
    google_pubsub_topic.additional-telemetry
  ]

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.additional-telemetry.id
    subfolder_matches = "test/directory"
  }

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.default-telemetry.id
    subfolder_matches = ""
  }

  state_notification_config {
    pubsub_topic_name = google_pubsub_topic.default-devicestatus.id
  }

  mqtt_config {
    mqtt_enabled_state = "MQTT_DISABLED"
  }

  http_config {
    http_enabled_state = "HTTP_DISABLED"
  }

  log_level = "INFO"

  credentials {
    public_key_certificate {
      format      = "X509_CERTIFICATE_PEM"
      certificate = file("test-fixtures/rsa_cert.pem")
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10), registryName, registryName)
}
