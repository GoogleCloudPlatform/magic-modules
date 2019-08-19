package google

import (
	"fmt"
	"reflect"
	"testing"
	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestValidateCloudIoTID(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foo"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a goog", Value: "googfoobar", ExpectError: true},
		{TestName: "starts with a number", Value: "1foobar", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has an backslash", Value: "foo\bar", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("f", 260), ExpectError: true},
	}

	es := testStringValidationCases(x, validateCloudIotID)
	if len(es) > 0 {
		t.Errorf("Failed to validate CloudIoT ID names: %v", es)
	}
}

func TestCloudIotRegistryStateUpgradeV0(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		V0State    map[string]interface{}
		V1Expected map[string]interface{}
	}{
		"Move single to plural": {
			V0State: map[string]interface{}{
				"event_notification_config": map[string]interface{}{
					"pubsub_topic_name": "projects/my-project/topics/my-topic",
				},
			},
			V1Expected: map[string]interface{}{
				"event_notification_configs": []interface{}{
					map[string]interface{}{
						"pubsub_topic_name": "projects/my-project/topics/my-topic",
					},
				},
			},
		},
		"Delete single if plural in state": {
			V0State: map[string]interface{}{
				"event_notification_config": map[string]interface{}{
					"pubsub_topic_name": "projects/my-project/topics/singular-topic",
				},
				"event_notification_configs": []interface{}{
					map[string]interface{}{
						"pubsub_topic_name": "projects/my-project/topics/plural-topic",
					},
				},
			},
			V1Expected: map[string]interface{}{
				"event_notification_configs": []interface{}{
					map[string]interface{}{
						"pubsub_topic_name": "projects/my-project/topics/plural-topic",
					},
				},
			},
		},
		"no-op": {
			V0State: map[string]interface{}{
				"name":      "my-test-name",
				"log_level": "INFO",
				"event_notification_configs": []interface{}{
					map[string]interface{}{
						"pubsub_topic_name": "projects/my-project/topics/plural-topic",
					},
				},
			},
			V1Expected: map[string]interface{}{
				"name":      "my-test-name",
				"log_level": "INFO",
				"event_notification_configs": []interface{}{
					map[string]interface{}{
						"pubsub_topic_name": "projects/my-project/topics/plural-topic",
					},
				},
			},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := resourceCloudIotRegistryStateUpgradeV0toV1(tc.V0State, &Config{})

			if err != nil {
				t.Error(err)
			}

			for k, v := range tc.V1Expected {
				if !reflect.DeepEqual(actual[k], v) {
					t.Errorf("expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
						k, v, k, actual[k], actual)
				}
			}
		})
	}
}

func TestAccCloudIoTRegistry_basic(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_basic(registryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudIoTRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudIoTRegistry_extended(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_extended(registryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudIoTRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudIoTRegistry_update(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_basic(registryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudIoTRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
			{
				Config: testAccCloudIoTRegistry_extended(registryName),
			},
			{
				Config: testAccCloudIoTRegistry_basic(registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudIoTRegistry_eventNotificationConfigDeprecatedSingleToPlural(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))
	topic := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				// Use deprecated field (event_notification_config) to create
				Config: testAccCloudIoTRegistry_singleEventNotificationConfig(topic, registryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_cloudiot_registry.foobar", "event_notification_configs.#", "1"),
				),
			},
			{
				// Use new field (event_notification_configs) to see if plan changed
				Config:             testAccCloudIoTRegistry_pluralEventNotificationConfigs(topic, registryName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccCloudIoTRegistry_eventNotificationConfigPluralToDeprecatedSingle(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))
	topic := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				// Use deprecated field (event_notification_config) to create
				Config: testAccCloudIoTRegistry_pluralEventNotificationConfigs(topic, registryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_cloudiot_registry.foobar", "event_notification_configs.#", "1"),
				),
			},
			{
				// Use new field (event_notification_configs) to see if plan changed
				Config:             testAccCloudIoTRegistry_singleEventNotificationConfig(topic, registryName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

func testAccCloudIoTRegistryExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientCloudIoT.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Registry does not exist")
		}
		return nil
	}
}

func testAccCloudIoTRegistry_basic(registryName string) string {
	return fmt.Sprintf(`
resource "google_cloudiot_registry" "foobar" {
	name = "%s"
}`, registryName)
}

func testAccCloudIoTRegistry_extended(registryName string) string {
	return fmt.Sprintf(`
resource "google_project_iam_binding" "cloud-iot-iam-binding" {
  members = ["serviceAccount:cloud-iot@system.gserviceaccount.com"]
  role    = "roles/pubsub.publisher"
}

resource "google_pubsub_topic" "default-devicestatus" {
  name = "psregistry-test-devicestatus-%s"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "psregistry-test-telemetry-%s"
}

resource "google_cloudiot_registry" "foobar" {
  depends_on = ["google_project_iam_binding.cloud-iot-iam-binding"]

  name = "%s"

  event_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.default-devicestatus.id}"
  }

  state_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.default-telemetry.id}"
  }

  http_config = {
    http_enabled_state = "HTTP_DISABLED"
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_DISABLED"
  }
	
  log_level = "INFO"

  credentials {
    public_key_certificate = {
      format      = "X509_CERTIFICATE_PEM"
      certificate = "${file("test-fixtures/rsa_cert.pem")}"
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10), registryName)
}

func testAccCloudIoTRegistry_singleEventNotificationConfig(topic, registryName string) string {
	return fmt.Sprintf(`
resource "google_project_iam_binding" "cloud-iot-iam-binding" {
  members = ["serviceAccount:cloud-iot@system.gserviceaccount.com"]
  role    = "roles/pubsub.publisher"
}

resource "google_pubsub_topic" "event-topic" {
  name = "%s"
}

resource "google_cloudiot_registry" "foobar" {
  depends_on = ["google_project_iam_binding.cloud-iot-iam-binding"]

  name = "%s"

  event_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.event-topic.id}"
  }
}
`, topic, registryName)
}

func testAccCloudIoTRegistry_pluralEventNotificationConfigs(topic, registryName string) string {
	return fmt.Sprintf(`
resource "google_project_iam_binding" "cloud-iot-iam-binding" {
  members = ["serviceAccount:cloud-iot@system.gserviceaccount.com"]
  role    = "roles/pubsub.publisher"
}

resource "google_pubsub_topic" "event-topic" {
  name = "%s"
}

resource "google_cloudiot_registry" "foobar" {
  depends_on = ["google_project_iam_binding.cloud-iot-iam-binding"]

  name = "%s"

  event_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.event-topic.id}"
  }
}
`, topic, registryName)
}
