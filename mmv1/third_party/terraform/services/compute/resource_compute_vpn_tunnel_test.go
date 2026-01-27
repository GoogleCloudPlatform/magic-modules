package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeVpnTunnel_regionFromGateway(t *testing.T) {
	t.Parallel()
	region := "us-central1"
	suffix := acctest.RandString(t, 10)
	if envvar.GetTestRegionFromEnv() == region {
		// Make sure we choose a region that isn't the provider default
		// in order to test getting the region from the gateway and not the
		// provider.
		region = "us-west1"
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeVpnTunnelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnel_regionFromGateway(suffix, region),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
		},
	})
}

func TestAccComputeVpnTunnel_router(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	router := fmt.Sprintf("tf-test-tunnel-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeVpnTunnelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnelRouter(suffix, router),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
		},
	})
}

func TestAccComputeVpnTunnel_routerWithSharedSecretWo_update(t *testing.T) {
	t.Parallel()

	router := fmt.Sprintf("tf-test-tunnel-%s", acctest.RandString(t, 10))
	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeVpnTunnelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnelRouterWithSharedSecretWo(suffix, router),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_vpn_tunnel.foobar", "shared_secret_wo"),
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.foobar", "shared_secret_wo_version", "1"),
				),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
			{
				Config: testAccComputeVpnTunnelRouterWithSharedSecretWo_update(suffix, router),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_vpn_tunnel.foobar", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_compute_vpn_tunnel.foobar", "shared_secret_wo"),
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.foobar", "shared_secret_wo_version", "2"),
				),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
		},
	})
}

func TestAccComputeVpnTunnel_defaultTrafficSelectors(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeVpnTunnelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnelDefaultTrafficSelectors(suffix),
			},
			{
				ResourceName:            "google_compute_vpn_tunnel.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shared_secret", "detailed_status"},
			},
		},
	})
}

// TestAccComputeVpnTunnel_cipherSuite tests the 'cipher_suite' block in the google_compute_vpn_tunnel resource.
func TestAccComputeVpnTunnel_cipherSuite(t *testing.T) {
	t.Parallel()

	// A unique name for the test resources
	suffix := acctest.RandString(t, 10)
	// Other necessary resources like network, gateway, etc. would be defined here.

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeVpnTunnelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test case 1: Basic cipher suite configuration
				Config: testAccComputeVpnTunnel_basicCipherSuite(suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.test_tunnel", "cipher_suite.0.phase1.0.encryption.0", "AES-GCM-16-128"),
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.test_tunnel", "cipher_suite.0.phase1.0.encryption.1", "AES-GCM-16-192"),
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.test_tunnel", "cipher_suite.0.phase2.0.integrity.0", "HMAC-SHA2-256-128"),
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.test_tunnel", "cipher_suite.0.phase2.0.integrity.1", "HMAC-SHA1-96"),
				),
			},
		},
	})
}

func TestAccComputeVpnTunnel_capacityTierZ2Z(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		CheckDestroy:             testAccCheckComputeVpnTunnelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVpnTunnel_capacityTierZ2Z(suffix, "DEFAULT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_vpn_tunnel.foobar", "capacity_tier", "DEFAULT"),
				),
			},
		},
	})
}

func testAccComputeVpnTunnel_capacityTierZ2Z(suffix, tier string) string {
	return fmt.Sprintf(`
provider "google-beta" {
  region = "us-east4"
}

data "google_project" "project" {
  provider = google-beta
}

resource "google_compute_network" "network" {
  provider                = google-beta
  name                    = "tf-test-network-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_interconnect" "interconnect" {
  provider             = google-beta
  name                 = "tf-test-interconnect-%[1]s"
  customer_name        = "internal_customer"
  interconnect_type    = "DEDICATED"
  link_type            = "LINK_TYPE_ETHERNET_100G_LR"
  location             = "https://www.googleapis.com/compute/v1/projects/${data.google_project.project.project_id}/global/interconnectLocations/z2z-us-east4-zone1-pniada-a"
  requested_link_count = 1
  admin_enabled        = true
}

resource "google_compute_router" "encrypted_router" {
  provider                      = google-beta
  name                          = "tf-test-encrypted-router-%[1]s"
  region                        = "us-east4"
  network                       = google_compute_network.network.id
  encrypted_interconnect_router = true
  bgp {
    asn = 64514
  }
}

resource "google_compute_router" "router" {
  provider = google-beta
  name     = "tf-test-router-%[1]s"
  region   = "us-east4"
  network  = google_compute_network.network.id
  bgp {
    asn = 64515
  }
}

resource "google_compute_interconnect_attachment" "attachment" {
  provider      = google-beta
  name          = "tf-test-attachment-%[1]s"
  interconnect  = google_compute_interconnect.interconnect.id
  type          = "DEDICATED"
  region        = "us-east4"
  bandwidth     = "BPS_10G"
  router        = google_compute_router.encrypted_router.id
  vlan_tag8021q = 1100
  encryption    = "IPSEC"
}

resource "google_compute_ha_vpn_gateway" "ha_gateway" {
  provider = google-beta
  name     = "tf-test-ha-gw-%[1]s"
  network  = google_compute_network.network.id
  region   = "us-east4"

  vpn_interfaces {
    id                      = 0
    interconnect_attachment = google_compute_interconnect_attachment.attachment.self_link
  }
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  provider        = google-beta
  name            = "tf-test-ext-gw-%[1]s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  provider                        = google-beta
  name                            = "tf-test-tunnel-%[1]s"
  region                          = "us-east4"
  vpn_gateway                     = google_compute_ha_vpn_gateway.ha_gateway.id
  vpn_gateway_interface           = 0
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0
  shared_secret                   = "unguessable"
  router                          = google_compute_router.router.id

  # Field under test
  capacity_tier                   = "%[2]s"
}
`, suffix, tier)
}

func testAccComputeVpnTunnel_basicCipherSuite(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-network-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-subnetwork-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "external-gateway-%[1]s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "tf-test-router-%[1]s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "test_tunnel" {
  name          = "tf-test-ha-vpn-tunnel-%[1]s"
  region        = "us-central1"
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.self_link
  vpn_gateway_interface           = 0 

  cipher_suite {
    phase1 {
      encryption = ["AES-GCM-16-128", "AES-GCM-16-192"]
    }
    phase2 {
      integrity  = ["HMAC-SHA2-256-128", "HMAC-SHA1-96"]
    }
  }
}
`, suffix)
}

func testAccComputeVpnTunnel_regionFromGateway(suffix, region string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "%[2]s"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "tf-test-%[1]s-esp"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "tf-test-%[1]s-udp500"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "tf-test-%[1]s-udp4500"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_vpn_tunnel" "foobar" {
  name                    = "tf-test-%[1]s"
  target_vpn_gateway      = google_compute_vpn_gateway.foobar.self_link
  shared_secret           = "unguessable"
  peer_ip                 = "8.8.8.8"
  local_traffic_selector  = [google_compute_subnetwork.foobar.ip_cidr_range]
  remote_traffic_selector = ["192.168.0.0/24", "192.168.1.0/24"]

  depends_on = [google_compute_forwarding_rule.foobar_udp4500]
}
`, suffix, region)
}

func testAccComputeVpnTunnelRouter(suffix, router string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-subnetwork-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "external-gateway-%[1]s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%[2]s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "tf-test-%[1]s"
  region             = google_compute_subnetwork.foobar.region
  vpn_gateway = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret      = "unguessable"
  router             = google_compute_router.foobar.self_link
  vpn_gateway_interface           = 0  
}
`, suffix, router)
}

func testAccComputeVpnTunnelRouterWithSharedSecretWo(suffix, router string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-subnetwork-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "external-gateway-%[1]s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%[2]s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name                            = "tf-test-%[1]s"
  region                          = google_compute_subnetwork.foobar.region
  vpn_gateway                     = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret_wo                = "I am write only, and should not be written to state"
  shared_secret_wo_version        = 1
  router                          = google_compute_router.foobar.self_link
  vpn_gateway_interface           = 0  
}
`, suffix, router)
}

func testAccComputeVpnTunnelRouterWithSharedSecretWo_update(suffix, router string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "foobar" {
  name          = "tf-test-subnetwork-%[1]s"
  network       = google_compute_network.foobar.self_link
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = google_compute_subnetwork.foobar.region
}

resource "google_compute_ha_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_subnetwork.foobar.region
}

resource "google_compute_external_vpn_gateway" "external_gateway" {
  name            = "external-gateway-%[1]s"
  redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
  description     = "An externally managed VPN gateway"
  interface {
    id         = 0
    ip_address = "8.8.8.8"
  }
}

resource "google_compute_router" "foobar" {
  name    = "%[2]s"
  region  = google_compute_subnetwork.foobar.region
  network = google_compute_network.foobar.self_link
  bgp {
    asn = 64514
  }
}

resource "google_compute_vpn_tunnel" "foobar" {
  name                            = "tf-test-%[1]s"
  region                          = google_compute_subnetwork.foobar.region
  vpn_gateway                     = google_compute_ha_vpn_gateway.foobar.id
  peer_external_gateway           = google_compute_external_vpn_gateway.external_gateway.id
  peer_external_gateway_interface = 0  
  shared_secret_wo                = "This is another secret, but still write only"
  shared_secret_wo_version        = 2
  router                          = google_compute_router.foobar.self_link
  vpn_gateway_interface           = 0  
}
`, suffix, router)
}

func testAccComputeVpnTunnelDefaultTrafficSelectors(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-%[1]s"
  auto_create_subnetworks = "true"
}

resource "google_compute_address" "foobar" {
  name   = "tf-test-%[1]s"
  region = "us-central1"
}

resource "google_compute_vpn_gateway" "foobar" {
  name    = "tf-test-%[1]s"
  network = google_compute_network.foobar.self_link
  region  = google_compute_address.foobar.region
}

resource "google_compute_forwarding_rule" "foobar_esp" {
  name        = "tf-test-%[1]s-esp"
  region      = google_compute_vpn_gateway.foobar.region
  ip_protocol = "ESP"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp500" {
  name        = "tf-test-%[1]s-udp500"
  region      = google_compute_forwarding_rule.foobar_esp.region
  ip_protocol = "UDP"
  port_range  = "500-500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_forwarding_rule" "foobar_udp4500" {
  name        = "tf-test-%[1]s-udp4500"
  region      = google_compute_forwarding_rule.foobar_udp500.region
  ip_protocol = "UDP"
  port_range  = "4500-4500"
  ip_address  = google_compute_address.foobar.address
  target      = google_compute_vpn_gateway.foobar.self_link
}

resource "google_compute_vpn_tunnel" "foobar" {
  name               = "tf-test-%[1]s"
  region             = google_compute_forwarding_rule.foobar_udp4500.region
  target_vpn_gateway = google_compute_vpn_gateway.foobar.self_link
  shared_secret      = "unguessable"
  peer_ip            = "8.8.8.8"
}
`, suffix)
}
