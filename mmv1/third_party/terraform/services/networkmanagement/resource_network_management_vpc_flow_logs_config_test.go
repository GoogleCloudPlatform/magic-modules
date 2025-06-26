package networkmanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkManagementVpcFlowLogsConfig_updateInterconnect(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_fullInterconnect(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.interconnect-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_updateInterconnect(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.interconnect-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_fullInterconnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "interconnect-test" {
  vpc_flow_logs_config_id = "tf-test-full-interconnect-test-id%{random_suffix}"
  location                = "global"
  interconnect_attachment = "projects/${data.google_project.project.number}/regions/us-east4/interconnectAttachments/${google_compute_interconnect_attachment.attachment.name}"
}

resource "google_compute_network" "network" {
  name     = "tf-test-full-interconnect-test-network%{random_suffix}"
}

resource "google_compute_router" "router" {
  name    = "tf-test-full-interconnect-test-router%{random_suffix}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "attachment" {
  name                     = "tf-test-full-interconnect-test-id%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
}

`, context)
}

func testAccNetworkManagementVpcFlowLogsConfig_updateInterconnect(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "interconnect-test" {
  vpc_flow_logs_config_id = "tf-test-full-interconnect-test-id%{random_suffix}"
  location                = "global"
  interconnect_attachment = "projects/${data.google_project.project.number}/regions/us-east4/interconnectAttachments/${google_compute_interconnect_attachment.attachment.name}"
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_30_SEC"
  description             = "This is an updated description"
  flow_sampling           = 0.5
  metadata                = "EXCLUDE_ALL_METADATA"
}

resource "google_compute_network" "network" {
  name     = "tf-test-full-interconnect-test-network%{random_suffix}"
}

resource "google_compute_router" "router" {
  name    = "tf-test-full-interconnect-test-router%{random_suffix}"
  network = google_compute_network.network.name
  bgp {
    asn = 16550
  }
}

resource "google_compute_interconnect_attachment" "attachment" {
  name                     = "tf-test-full-interconnect-test-id%{random_suffix}"
  edge_availability_domain = "AVAILABILITY_DOMAIN_1"
  type                     = "PARTNER"
  router                   = google_compute_router.router.id
  mtu                      = 1500
}

`, context)
}

func TestAccNetworkManagementVpcFlowLogsConfig_updateVpn(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_fullVpn(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_updateVpn(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_fullVpn(context map[string]interface{}) string {
	vpcFlowLogsCfg := acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "example" {
  vpc_flow_logs_config_id = "id-example-%{random_suffix}"
  location                = "global"
  vpn_tunnel              = "projects/${data.google_project.project.number}/regions/us-central1/vpnTunnels/${google_compute_vpn_tunnel.tunnel.name}"
}
`, context)
	return fmt.Sprintf("%s\n\n%s\n\n", vpcFlowLogsCfg, testAccNetworkManagementVpcFlowLogsConfig_baseResources(context))
}

func testAccNetworkManagementVpcFlowLogsConfig_updateVpn(context map[string]interface{}) string {
	vpcFlowLogsCfg := acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_network_management_vpc_flow_logs_config" "example" {
  vpc_flow_logs_config_id = "id-example-%{random_suffix}"
  location                = "global"
  vpn_tunnel              = "projects/${data.google_project.project.number}/regions/us-central1/vpnTunnels/${google_compute_vpn_tunnel.tunnel.name}"
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_30_SEC"
  description             = "This is an updated description"
  flow_sampling           = 0.5
  metadata                = "EXCLUDE_ALL_METADATA"
}
`, context)
	return fmt.Sprintf("%s\n\n%s\n\n", vpcFlowLogsCfg, testAccNetworkManagementVpcFlowLogsConfig_baseResources(context))
}

func TestAccNetworkManagementVpcFlowLogsConfig_network(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_network(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.network-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_networkUpdate(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.network-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_network(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name = "tf-test-flow-logs-network-%{random_suffix}"
}

resource "google_network_management_vpc_flow_logs_config" "network-test" {
  vpc_flow_logs_config_id = "tf-test-network-id-%{random_suffix}"
  location                = "global"
  network                 = google_compute_network.network.id
  state                   = "ENABLED"
}
`, context)
}

func testAccNetworkManagementVpcFlowLogsConfig_networkUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name = "tf-test-flow-logs-network-%{random_suffix}"
}

resource "google_network_management_vpc_flow_logs_config" "network-test" {
  vpc_flow_logs_config_id = "tf-test-network-id-%{random_suffix}"
  location                = "global"
  network                 = google_compute_network.network.id

  // Updated fields
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_10_MIN"
  flow_sampling           = 0.05
  metadata                = "INCLUDE_ALL_METADATA"
  description             = "Updated description for network test"
}
`, context)
}

func TestAccNetworkManagementVpcFlowLogsConfig_subnet(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_subnet(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.subnet-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_subnetUpdate(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.subnet-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_subnet(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-subnet-network-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-flow-logs-subnet-%{random_suffix}"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.network.id
}

resource "google_network_management_vpc_flow_logs_config" "subnet-test" {
  vpc_flow_logs_config_id = "tf-test-subnet-id-%{random_suffix}"
  location                = "global"
  subnet                  = google_compute_subnetwork.subnet.id
}
`, context)
}

func testAccNetworkManagementVpcFlowLogsConfig_subnetUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network" {
  name                    = "tf-test-subnet-network-%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "tf-test-flow-logs-subnet-%{random_suffix}"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.network.id
}

resource "google_network_management_vpc_flow_logs_config" "subnet-test" {
  vpc_flow_logs_config_id = "tf-test-subnet-id-%{random_suffix}"
  location                = "global"
  subnet                  = google_compute_subnetwork.subnet.id
  
  // Updated fields
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_30_SEC"
  flow_sampling           = 0.5
  metadata                = "EXCLUDE_ALL_METADATA"
  description             = "Updated description for subnet test"
}
`, context)
}

func testAccNetworkManagementVpcFlowLogsConfig_baseResources(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_vpn_tunnel" "tunnel" {
  name               = "tf-test-example-tunnel%{random_suffix}"
  peer_ip            = "15.0.0.120"
  shared_secret      = "a secret message"
  target_vpn_gateway = google_compute_vpn_gateway.target_gateway.id

  depends_on = [
    google_compute_forwarding_rule.fr_esp,
    google_compute_forwarding_rule.fr_udp500,
    google_compute_forwarding_rule.fr_udp4500,
  ]
}

resource "google_compute_vpn_gateway" "target_gateway" {
  name     = "tf-test-example-gateway%{random_suffix}"
  network  = google_compute_network.network.id
}

resource "google_compute_network" "network" {
  name     = "tf-test-example-network%{random_suffix}"
}

resource "google_compute_address" "vpn_static_ip" {
  name     = "tf-test-example-address%{random_suffix}"
}

resource "google_compute_forwarding_rule" "fr_esp" {
  name        = "tf-test-example-fresp%{random_suffix}"
  ip_protocol = "ESP"
  ip_address  = google_compute_address.vpn_static_ip.address
  target      = google_compute_vpn_gateway.target_gateway.id
}

resource "google_compute_forwarding_rule" "fr_udp500" {
  name        = "tf-test-example-fr500%{random_suffix}"
  ip_protocol = "UDP"
  port_range  = "500"
  ip_address  = google_compute_address.vpn_static_ip.address
  target      = google_compute_vpn_gateway.target_gateway.id
}

resource "google_compute_forwarding_rule" "fr_udp4500" {
  name        = "tf-test-example-fr4500%{random_suffix}"
  ip_protocol = "UDP"
  port_range  = "4500"
  ip_address  = google_compute_address.vpn_static_ip.address
  target      = google_compute_vpn_gateway.target_gateway.id
}

resource "google_compute_route" "route" {
  name                = "tf-test-example-route%{random_suffix}"
  network             = google_compute_network.network.name
  dest_range          = "15.0.0.0/24"
  priority            = 1000
  next_hop_vpn_tunnel = google_compute_vpn_tunnel.tunnel.id
}
`, context)
}

func TestAccNetworkManagementVpcFlowLogsConfig_organization(t *testing.T) {
	t.Parallel()
	// This helper function retrieves the organization ID from the test environment.
	// For local runs, it reads the GOOGLE_ORG environment variable.
	// In CI, the test runner provides the correct test organization ID.
	orgID := acctest.GetOrgID(t)

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"org_id":        orgID,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkManagementVpcFlowLogsConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_organization_basic_with_iam(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.org-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "vpc_flow_logs_config_id"},
			},
			{
				Config: testAccNetworkManagementVpcFlowLogsConfig_organization_update_with_iam(context),
			},
			{
				ResourceName:            "google_network_management_vpc_flow_logs_config.org-test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "vpc_flow_logs_config_id"},
			},
		},
	})
}

func testAccNetworkManagementVpcFlowLogsConfig_organization_basic_with_iam(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

data "google_service_account" "test_runner" {
    account_id = "hashicorp-test-runner"
    project    = data.google_project.project.project_id
}

// This resource grants the 'Network Management Admin' role to the test service account.
resource "google_organization_iam_member" "nm_admin_binding" {
    org_id = "%{org_id}"
    role   = "roles/networkmanagement.admin"
    member = "serviceAccount:${data.google_service_account.test_runner.email}"
}

resource "google_network_management_vpc_flow_logs_config" "org-test" {
  vpc_flow_logs_config_id = "tf-test-org-id-%{random_suffix}"
  location                = "global"
  parent                  = "organizations/%{org_id}"
  state                   = "ENABLED"

  // This ensures the IAM role is granted before this resource is created.
  depends_on = [google_organization_iam_member.nm_admin_binding]
}
`, context)
}

func testAccNetworkManagementVpcFlowLogsConfig_organization_update_with_iam(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

data "google_service_account" "test_runner" {
    account_id = "hashicorp-test-runner"
    project    = data.google_project.project.project_id
}

resource "google_organization_iam_member" "nm_admin_binding" {
    org_id = "%{org_id}"
    role   = "roles/networkmanagement.admin"
    member = "serviceAccount:${data.google_service_account.test_runner.email}"
}

resource "google_network_management_vpc_flow_logs_config" "org-test" {
  vpc_flow_logs_config_id = "tf-test-org-id-%{random_suffix}"
  location                = "global"
  parent                  = "organizations/%{org_id}"

  // Updated fields
  state                   = "DISABLED"
  aggregation_interval    = "INTERVAL_15_MIN"
  flow_sampling           = 0.2
  metadata                = "INCLUDE_ALL_METADATA"
  description             = "Updated description for org test"
  
  depends_on = [google_organization_iam_member.nm_admin_binding]
}
`, context)
}
