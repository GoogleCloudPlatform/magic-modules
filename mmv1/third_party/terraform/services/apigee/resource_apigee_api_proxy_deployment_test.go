package apigee_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccApigeeApiProxyDeployment_apigeeApiProxyDeploymentBasicExample(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()
	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeApiProxyDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeApiProxyDeployment_apigeeApiProxyDeploymentBasicExample(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_apigee_api_proxy_deployment.deploy", "state"),
				),
			},
			{
				ResourceName:            "google_apigee_api_proxy_deployment.deploy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"override", "sequenced_rollout", "service_account"},
			},
		},
	})
}

func testAccApigeeApiProxyDeployment_apigeeApiProxyDeploymentBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "project" {
		project_id      = "tf-test%{random_suffix}"
		name            = "tf-test%{random_suffix}"
		org_id          = "%{org_id}"
		billing_account = "%{billing_account}"
		deletion_policy = "DELETE"
	}

	resource "google_project_service" "apigee" {
		project = google_project.project.project_id
		service = "apigee.googleapis.com"
	}

	resource "google_project_service" "servicenetworking" {
		project    = google_project.project.project_id
		service    = "servicenetworking.googleapis.com"
		depends_on = [google_project_service.apigee]
	}
		
	resource "google_project_service" "compute" {
		project    = google_project.project.project_id
		service    = "compute.googleapis.com"
		depends_on = [google_project_service.servicenetworking]
	}

	resource "time_sleep" "wait_120_seconds" {
		create_duration = "120s"
		depends_on      = [google_project_service.compute]
	}
		
	resource "google_compute_network" "apigee_network" {
		name       = "apigee-network%{random_suffix}"
		project    = google_project.project.project_id
		depends_on = [time_sleep.wait_120_seconds]
	}
	
	resource "google_compute_global_address" "apigee_range" {
		name          = "tf-test-apigee-range%{random_suffix}"
		purpose       = "VPC_PEERING"
		address_type  = "INTERNAL"
		prefix_length = 16
		network       = google_compute_network.apigee_network.id
		project       = google_project.project.project_id
	}
		
	resource "google_service_networking_connection" "apigee_vpc_connection" {
		network                 = google_compute_network.apigee_network.id
		service                 = "servicenetworking.googleapis.com"
		reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
	}
		
	resource "google_apigee_organization" "apigee_org" {
		analytics_region   = "us-central1"
		project_id         = google_project.project.project_id
		authorized_network = google_compute_network.apigee_network.id
		depends_on         = [google_service_networking_connection.apigee_vpc_connection]
	}

	resource "google_apigee_environment" "env" {
		org_id       = google_apigee_organization.apigee_org.id
		name         = "dev"
		display_name = "dev"
		description  = "terraform test env"

	}

	resource "google_apigee_instance" "apigee_ins" {
		name         = "apigee-instance%{random_suffix}"
		location     = "us-central1"
		org_id       = google_apigee_organization.apigee_org.id
		depends_on   = [google_apigee_environment.env]
	}
		
	resource "google_apigee_instance_attachment" "instance_att" {
		instance_id  = google_apigee_instance.apigee_ins.id
		environment  = google_apigee_environment.env.name
		depends_on   = [google_apigee_instance.apigee_ins]
	}
		
	resource "google_apigee_api" "proxy" {
		name          = "tf-test-apigee-proxy"
		org_id        = google_project.project.project_id
		config_bundle = "./test-fixtures/apigee_api_bundle.zip"
		depends_on    = [google_apigee_instance_attachment.instance_att]
	}

	resource "google_apigee_api_proxy_deployment" "deploy" {
		org               = google_project.project.project_id
		environment       = google_apigee_environment.env.name
		api               = google_apigee_api.proxy.name
		revision          = 1
		override          = true
		sequenced_rollout = true
	}
`, context)
}

func testAccCheckApigeeApiProxyDeploymentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_apigee_api_proxy_deployment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs,
				"{{ApigeeBasePath}}organizations/{{org}}/environments/{{environment}}/apis/{{api}}/revisions/{{revision}}/deployments")
			if err != nil {
				return err
			}
			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("Apigee API proxy revision still appears deployed at %s", url)
			}
		}
		return nil
	}
}
