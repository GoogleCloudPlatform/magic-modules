package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkConnectivityServiceConnectionPolicy_update(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-network-%s", RandString(t, 10))
	networkProducerName := fmt.Sprintf("tf-test-network-%s", RandString(t, 10))
	subnetworkConsumerName := fmt.Sprintf("tf-test-subnet-consumer-%s", RandString(t, 10))
	subnetworkProducerName := fmt.Sprintf("tf-test-subnet-producer-%s", RandString(t, 10))
	serviceConnectionPolicyName := fmt.Sprintf("tf-test-service-connection-policy-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityServiceConnectionPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityServiceConnectionPolicy_basic(networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName),
			},
			{
				ResourceName:      "google_network_connectivity_service_connection_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkConnectivityServiceConnectionPolicy_update(networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName),
			},
			{
				ResourceName:      "google_network_connectivity_service_connection_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkConnectivityServiceConnectionPolicy_basic(networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName),
			},
			{
				ResourceName:      "google_network_connectivity_service_connection_policy.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkConnectivityServiceConnectionPolicy_basic(networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName string) string {
	return fmt.Sprintf(`
  resource "google_compute_network" "consumer_net" {
    name                    = "%s"
    auto_create_subnetworks = false
  }
  
  resource "google_compute_subnetwork" "consumer_subnet" {
    name          = "%s"
    ip_cidr_range = "10.0.0.0/16"
    region        = "us-central1"
    network       = google_compute_network.consumer_net.id
  }
  
  resource "google_compute_network" "producer_net" {
    name                    = "%s"
    auto_create_subnetworks = false
  }
  
  resource "google_compute_subnetwork" "producer_subnet" {
    name          = "%s"
    ip_cidr_range = "10.0.0.0/16"
    region        = "us-central1"
    network       = google_compute_network.producer_net.id
  }
  
  resource "google_network_connectivity_service_connection_policy" "default" {
    name = "%s"
    location = "us-central1"
    description = "my basic sevice connection policy"
    service_class = "gcp-memorystore-redis"
    network = google_compute_network.producer_net.id
    psc_config {
      subnetworks = [google_compute_subnetwork.producer_subnet.id]
      limit = 2
    }
  }
`, networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName)
}

func testAccNetworkConnectivityServiceConnectionPolicy_update(networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "consumer_net" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "consumer_subnet" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.consumer_net.id
}

resource "google_compute_network" "producer_net" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "producer_subnet" {
  name          = "%s"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.producer_net.id
}

resource "google_network_connectivity_service_connection_policy" "default" {
  name = "%s"
  location = "us-central1"
  description = "my basic sevice connection policy"
  service_class = "gcp-memorystore-redis"
  network = google_compute_network.producer_net.id
  psc_config {
    subnetworks = [google_compute_subnetwork.producer_subnet.id]
    limit = 2
  }
  labels      = {
    foo = "bar"
  }
}
`, networkName, subnetworkConsumerName, networkProducerName, subnetworkProducerName, serviceConnectionPolicyName)
}
