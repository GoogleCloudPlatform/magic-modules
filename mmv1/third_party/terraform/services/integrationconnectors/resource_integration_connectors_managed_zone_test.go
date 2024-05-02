package integrationconnectors_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckIntegrationConnectorsManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_full(context),
			},
			{
				ResourceName:            "google_integration_connectors_managed_zone.samplemanagedzone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
			{
				Config: testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_update(context),
			},
			{
				ResourceName:            "google_integration_connectors_managed_zone.samplemanagedzone",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_integration_connectors_managed_zone" "samplemanagedzone" {
  name     = "tf-test-test-managed-zone%{random_suffix}"
  description = "tf created description"
  labels = {
    intent = "example"
  }
  target_project="connectors-example"
  target_vpc="default"
  dns="connectors.example.com."
}
`, context)
}

func testAccIntegrationConnectorsManagedZone_integrationConnectorsManagedZoneExample_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "test_project" {
}

resource "google_integration_connectors_managed_zone" "samplemanagedzone" {
  name     = "tf-test-test-managed-zone%{random_suffix}"
  description = "tf created description"
  labels = {
    example = "sample"
  }
  target_project="connectors-ip-test"
  target_vpc="cp-ga-bug-bash"
  dns="connectors-new.example.com."
}
`, context)
}
