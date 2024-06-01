package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/client"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/securitycenter/v1"
)

func TestAccSecurityCenterOrganizationNotificationConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"topic_name":    acctest.RandString(t, 10),
		"config_id":     acctest.RandString(t, 10),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterOrganizationNotificationConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterOrganizationNotificationConfig_basic(context),
			},
			{
				ResourceName:            "google_scc_organization_notification_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "config_id"},
			},
		},
	})
}

func testAccSecurityCenterOrganizationNotificationConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_pubsub_topic" "scc_v2_organization_notification_config" {
  name = "tf-test-topic-%{random_suffix}"
}

resource "google_scc_organization_notification_config" "default" {
  config_id    = "tf-test-config-%{random_suffix}"
  organization = "%{org_id}"
  location     = "global"
  description  = "A test organization notification config"
  pubsub_topic = google_pubsub_topic.scc_v2_organization_notification_config.id

  streaming_config {
    filter = "severity = \"HIGH\""
  }
}
`, context)
}

func testAccCheckSecurityCenterOrganizationNotificationConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.Provider.Meta().(*client.Config)
		sccService, err := securitycenter.NewService(context.Background(), config.GoogleClientOptions...)
		if err != nil {
			return fmt.Errorf("Error creating Security Command Center client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_organization_notification_config" {
				continue
			}

			orgID := rs.Primary.Attributes["organization"]
			configID := rs.Primary.Attributes["config_id"]
			location := rs.Primary.Attributes["location"]
			name := fmt.Sprintf("organizations/%s/locations/%s/notificationConfigs/%s", orgID, location, configID)

			_, err := sccService.Organizations.Locations.NotificationConfigs.Get(name).Do()
			if err == nil {
				return fmt.Errorf("Notification config %s still exists", name)
			}
			if !isGoogleAPINotFoundError(err) {
				return fmt.Errorf("Error checking if Notification config %s exists: %s", name, err)
			}
		}
		return nil
	}
}

func isGoogleAPINotFoundError(err error) bool {
	if err == nil {
		return false
	}
	apiErr, ok := err.(*googleapi.Error)
	if !ok {
		return false
	}
	return apiErr.Code == 404
}



