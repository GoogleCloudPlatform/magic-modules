package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServiceNetworkingConnection_create(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-%s", randString(t, 10))
	addr := fmt.Sprintf("tf-test-%s", randString(t, 10))
	service := "servicenetworking.googleapis.com"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testServiceNetworkingConnectionDestroy(t, service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(network, addr, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceNetworkingConnection_update(t *testing.T) {
	t.Parallel()

	network := fmt.Sprintf("tf-test-%s", randString(t, 10))
	addr1 := fmt.Sprintf("tf-test-%s", randString(t, 10))
	addr2 := fmt.Sprintf("tf-test-%s", randString(t, 10))
	service := "servicenetworking.googleapis.com"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testServiceNetworkingConnectionDestroy(t, service, network),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceNetworkingConnection(network, addr1, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      testAccServiceNetworkingConnection_newAddr(network, addr1, addr2, "servicenetworking.googleapis.com"),
				ExpectError: regexp.MustCompile("%*Cannot modify allocated ranges in CreateConnection%*"),
			},
			{
				Config: testAccServiceNetworkingConnection_appendAddr(network, addr1, addr2, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceNetworkingConnection(network, addr1, "servicenetworking.googleapis.com"),
			},
			{
				ResourceName:      "google_service_networking_connection.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testServiceNetworkingConnectionDestroy(t *testing.T, parent, network string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		parentService := "services/" + parent
		networkName := fmt.Sprintf("projects/%s/global/networks/%s", getTestProjectFromEnv(), network)
		listCall := config.NewServiceNetworkingClient(config.userAgent).Services.Connections.List(parentService).Network(networkName)
		if config.UserProjectOverride {
			listCall.Header().Add("X-Goog-User-Project", getTestProjectFromEnv())
		}
		response, err := listCall.Do()
		if err != nil {
			return err
		}

		for _, c := range response.Connections {
			if c.Network == networkName {
				return fmt.Errorf("Found %s which should have been destroyed.", networkName)
			}
		}

		return nil
	}
}

func testAccServiceNetworkingConnection(networkName, addressRangeName, serviceName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "servicenet" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = google_compute_network.servicenet.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar.name]
}
`, networkName, addressRangeName, serviceName)
}

// this config is in addition to the config above this, it adds a new
// address and tries to create a new service networking connection with
// the same service and new address - this should cause an error
func testAccServiceNetworkingConnection_newAddr(networkName, addr1, addr2, serviceName string) string {
	return fmt.Sprintf(`
%s

resource "google_compute_global_address" "foobar_two" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar_two" {
  network                 = google_compute_network.servicenet.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar_two.name]
}
`, testAccServiceNetworkingConnection(networkName, addr1, serviceName), addr2, serviceName)
}

// this configuration will keep the new address created, but alter the original
// service connection (instead of creating a new one) and rather append the new
// address to the reserved_peering_ranges
func testAccServiceNetworkingConnection_appendAddr(networkName, addr1, addr2, serviceName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "servicenet" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_global_address" "foobar" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.servicenet.self_link
}

resource "google_compute_global_address" "foobar_two" {
  name          = "%s"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.servicenet.self_link
}

resource "google_service_networking_connection" "foobar" {
  network                 = google_compute_network.servicenet.self_link
  service                 = "%s"
  reserved_peering_ranges = [google_compute_global_address.foobar.name, google_compute_global_address.foobar_two.name]
}
`, networkName, addr1, addr2, serviceName)
}
