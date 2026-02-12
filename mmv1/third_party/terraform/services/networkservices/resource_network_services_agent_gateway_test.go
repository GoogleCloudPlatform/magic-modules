package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdate(t *testing.T) {
	t.Skip("b/484137930")
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesAgentGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdate(context),
			},
			{
				ResourceName:            "google_network_services_agent_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayFullExample(context),
			},
			{
				ResourceName:            "google_network_services_agent_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdate(context),
			},
			{
				ResourceName:            "google_network_services_agent_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_agent_gateway" "default" {
  name     = "tf-test-my-full-agent-gateway%{random_suffix}"
  location = "us-central1"
  description = "A very full configuration for Agent Gateway"
  labels = {
    env = "prod"
    tier = "silver"
  }

  protocols = []
  google_managed {
    governed_access_path = "CLIENT_TO_AGENT"
  }

  registries = [
    "//agentregistry.googleapis.com/projects/%{project}/locations/us-central1"
  ]

  network_config {
    egress {
      network_attachment = "projects/%{project}/regions/us-central1/networkAttachments/oh-my-network-attachment"
    }
  }
}
`, context)
}
