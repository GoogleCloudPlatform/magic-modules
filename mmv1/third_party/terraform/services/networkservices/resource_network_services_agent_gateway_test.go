package networkservices_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/networkservices"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdate(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"name":          "tf-test-my-full-agent-gateway" + randomSuffix,
		"random_suffix": randomSuffix,
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
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdateStep3(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_network_services_agent_gateway.default",
							plancheck.ResourceActionUpdate,
						),
					},
				},
			},
			{
				ResourceName:            "google_network_services_agent_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_network_services_agent_gateway.default",
							plancheck.ResourceActionUpdate,
						),
					},
				},
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
  name     = "%{name}"
  location = "us-central1"
  description = "A very full configuration for Agent Gateway"
  labels = {
    env = "prod"
    tier = "silver"
  }

  protocols = ["MCP"]
  google_managed {
    governed_access_path = "AGENT_TO_ANYWHERE"
  }

  registries = [
    "//agentregistry.googleapis.com/projects/%{project}/locations/us-central1"
  ]

  network_config {
    egress {
      network_attachment = google_compute_network_attachment.default.id
    }
  }

  depends_on = [google_project_service.agent_registry]
}
`, context) + testAccNetworkServicesAgentGateway_sharedNetworkResources(context)
}

func testAccNetworkServicesAgentGateway_sharedNetworkResources(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project_service" "agent_registry" {
  service            = "agentregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_compute_network" "default" {
  name                    = "net-%{name}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "subnet-%{name}"
  region        = "us-central1"
  network       = google_compute_network.default.id
  ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_network_attachment" "default" {
  name                  = "na-%{name}"
  region                = "us-central1"
  connection_preference = "ACCEPT_AUTOMATIC"
  subnetworks           = [google_compute_subnetwork.default.self_link]
}
`, context)
}

func testAccNetworkServicesAgentGateway_networkServicesAgentGatewayUpdateStep3(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_services_agent_gateway" "default" {
  name     = "%{name}"
  location = "us-central1"
  description = "A full configuration for Agent Gateway"
  labels = {
    env = "prod"
    tier = "silver"
  }

  protocols = ["MCP"]
  google_managed {
    governed_access_path = "AGENT_TO_ANYWHERE"
  }

  registries = [
    "//agentregistry.googleapis.com/projects/%{project}/locations/us-central1"
  ]

  network_config {
    egress {
      network_attachment = google_compute_network_attachment.default.id
    }
  }

  depends_on = [google_project_service.agent_registry]
}
`, context) + testAccNetworkServicesAgentGateway_sharedNetworkResources(context)
}

func TestAccNetworkServicesAgentGateway_networkServicesAgentGatewayAgentConnectivityTemplateUpdate(t *testing.T) {
	t.Parallel()

	acctest.SkipTestUntil(t, "2026-10-31") // b/548696678: Backend connectivity actuation rollout pending

	randomSuffix := acctest.RandString(t, 10)
	template1 := bootstrapAgentConnectivityTemplate(t, "tf-test-act-1-"+randomSuffix, "us-central1")
	template2 := bootstrapAgentConnectivityTemplate(t, "tf-test-act-2-"+randomSuffix, "us-central1")

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"name":          "tf-test-my-agent-gateway" + randomSuffix,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesAgentGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayAgentConnectivityTemplateUpdate(context, template1),
			},
			{
				ResourceName:            "google_network_services_agent_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesAgentGateway_networkServicesAgentGatewayAgentConnectivityTemplateUpdate(context, template2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"google_network_services_agent_gateway.default",
							plancheck.ResourceActionUpdate,
						),
					},
				},
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

func testAccNetworkServicesAgentGateway_networkServicesAgentGatewayAgentConnectivityTemplateUpdate(context map[string]interface{}, template string) string {
	ctx := make(map[string]interface{})
	for k, v := range context {
		ctx[k] = v
	}
	ctx["template"] = template
	return acctest.Nprintf(`
resource "google_network_services_agent_gateway" "default" {
  name     = "%{name}"
  location = "us-central1"

  google_managed {
    governed_access_path = "CLIENT_TO_AGENT"
  }

  agent_connectivity_template = "%{template}"

  registries = [
    "//agentregistry.googleapis.com/projects/%{project}/locations/us-central1"
  ]

  depends_on = [google_project_service.agent_registry]
}

resource "google_project_service" "agent_registry" {
  service            = "agentregistry.googleapis.com"
  disable_on_destroy = false
}
`, ctx)
}

func bootstrapAgentConnectivityTemplate(t *testing.T, templateId, location string) string {
	t.Helper()
	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap config.")
	}

	project := envvar.GetTestProjectFromEnv()
	baseURL := transport_tpg.BaseUrl(networkservices.Product, config)
	getURL := fmt.Sprintf("%sprojects/%s/locations/%s/agentConnectivityTemplates/%s",
		baseURL, project, location, templateId)

	headers := make(http.Header)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Headers:   headers,
	})

	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		postURL := fmt.Sprintf("%sprojects/%s/locations/%s/agentConnectivityTemplates?agentConnectivityTemplateId=%s",
			baseURL, project, location, templateId)
		obj := map[string]interface{}{
			"description": "Bootstrapped AgentConnectivityTemplate for acceptance test",
		}
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    postURL,
			UserAgent: config.UserAgent,
			Body:      obj,
			Headers:   headers,
		})
		if err != nil {
			t.Fatalf("Error creating bootstrapped AgentConnectivityTemplate %s: %s", templateId, err)
		}

		err = networkservices.NetworkServicesOperationWaitTime(
			config, res, project, "Creating AgentConnectivityTemplate", config.UserAgent,
			4*time.Minute)
		if err != nil {
			t.Fatalf("Error waiting for bootstrapped AgentConnectivityTemplate %s: %s", templateId, err)
		}
	}

	t.Cleanup(func() {
		deleteURL := fmt.Sprintf("%sprojects/%s/locations/%s/agentConnectivityTemplates/%s",
			baseURL, project, location, templateId)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "DELETE",
			Project:   project,
			RawURL:    deleteURL,
			UserAgent: config.UserAgent,
		})
		if err == nil {
			_ = networkservices.NetworkServicesOperationWaitTime(
				config, res, project, "Deleting AgentConnectivityTemplate", config.UserAgent,
				4*time.Minute)
		}
	})

	return fmt.Sprintf("projects/%s/locations/%s/agentConnectivityTemplates/%s", project, location, templateId)
}
