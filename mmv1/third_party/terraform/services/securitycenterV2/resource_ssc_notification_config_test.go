package securitycenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterOrganizationNotificationConfig_updateStreamingConfigFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"location":      "global",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterNotificationConfigV2DestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterOrganizationNotificationConfig_sccNotificationConfigBasicExample(context),
			},
			{
				ResourceName:            "google_scc_organization_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization", "config_id"},
			},
			{
				Config: testAccSecurityCenterOrganizationNotificationConfig_updateStreamingConfigFilter(context),
			},
			{
				ResourceName:            "google_scc_organization_notification_config.custom_notification_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization", "config_id"},
			},
		},
	})
}

func testAccSecurityCenterOrganizationNotificationConfig_updateStreamingConfigFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_organization_notification" {
  name = "tf-test-my-topic%{random_suffix}"
}

resource "google_scc_organization_notification_config" "custom_organization_notification_config" {
  config_id    = "tf-test-my-config%{random_suffix}"
  organization = "%{org_id}"
  location     = "%{location}"
  description  = "My custom Cloud Security Command Center Finding Notification Configuration"
  pubsub_topic =  google_pubsub_topic.scc_organization_notification.id

  streaming_config {
    filter = "category = \"OPEN_FIREWALL\""
  }
}
`, context)
}

func testAccSecurityCenterOrganizationNotificationConfig_sccNotificationConfigBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_organization_notification" {
  name = "tf-test-my-topic%{random_suffix}"
}

resource "google_scc_organization_notification_config" "custom_organization_notification_config" {
  config_id    = "tf-test-my-config%{random_suffix}"
  organization = "%{org_id}"
  location     = "%{location}"
  description  = "My custom Cloud Security Command Center Finding Notification Configuration"
  pubsub_topic =  google_pubsub_topic.scc_organization_notification.id

  streaming_config {
    filter = "category = \"FIREWALL_EVENT\""
  }
}
`, context)
}

func testAccCheckSecurityCenterOrganizationNotificationConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *resource.State) error {
		configs := getNotificationConfigResources(s)

		for _, config := range configs {
			organization := config.Primary.Attributes["organization"]
			location := config.Primary.Attributes["location"]
			configID := config.Primary.Attributes["config_id"]

			// Construct the resource name based on the attributes
			resourceName := fmt.Sprintf("organizations/%s/locations/%s/notificationConfigs/%s", organization, location, configID)

			// Initialize the Google API client
			config := acctest.Provider.Meta().(*Config)
			sccService, err := securitycenter.NewService(config.client)
			if err != nil {
				return fmt.Errorf("Error creating Security Command Center service: %s", err)
			}

			// Check if the Notification Config still exists
			_, err = sccService.Organizations.Locations.NotificationConfigs.Get(resourceName).Do()
			if err == nil {
				return fmt.Errorf("Notification Config %s still exists", resourceName)
			}
			if !isNotFoundError(err) {
				return fmt.Errorf("Error checking for Notification Config %s: %s", resourceName, err)
			}
		}

		return nil
	}
}

// Helper function to get the relevant resources from the state
func getNotificationConfigResources(s *resource.State) []*resource.Resource {
	var resources []*resource.Resource
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "google_scc_organization_notification_config" {
			resources = append(resources, rs)
		}
	}
	return resources
}

// Helper function to check if the error is a NotFound error
func isNotFoundError(err error) bool {
	apiErr, ok := err.(*googleapi.Error)
	return ok && apiErr.Code == 404
}
