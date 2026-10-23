package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkConnectivityServiceConnectionMap_networkConnectivityMapBasicUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_id":         envvar.GetTestProjectFromEnv(),
		"service_class_name": "gcp-memorystore-redis",
		"random_suffix":      acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityServiceConnectionMapDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityServiceConnectionMap_networkConnectivityMapBasicCreate(context),
			},
			{
				ResourceName:            "google_network_connectivity_service_connection_map.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivityServiceConnectionMap_networkConnectivityMapBasicUpdate(context),
			},
			{
				ResourceName:            "google_network_connectivity_service_connection_map.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkConnectivityServiceConnectionMap_networkConnectivityMapBasicCreate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_service_connection_map" "default" {
  name = "tf-test-my-network-connectivity-map%{random_suffix}"
  description   = "my basic service connection map"
  location = "us-central1"
  service_class = "%{service_class_name}"
  producer_psc_configs {
    service_attachment_uri = google_compute_service_attachment.psc_ilb_service_attachment.id
  }
  consumer_psc_configs {
    network = google_compute_network.consumer_net.id
    project = "%{project_id}"
    consumer_instance_project = "%{project_id}"
  }
}

resource "google_compute_service_attachment" "psc_ilb_service_attachment" {
  name        = "tf-test-my-psc-ilb%{random_suffix}"
  region      = "us-central1"
  description = "A service attachment configured with Terraform"

  enable_proxy_protocol    = true
  connection_preference    = "ACCEPT_AUTOMATIC"
  nat_subnets              = [google_compute_subnetwork.psc_ilb_nat.id]
  target_service           = google_compute_forwarding_rule.psc_ilb_target_service.id
}

resource "google_compute_address" "psc_ilb_consumer_address" {
  name   = "tf-test-psc-ilb-consumer-address%{random_suffix}"
  region = "us-central1"

  subnetwork   = "default"
  address_type = "INTERNAL"
}

resource "google_compute_forwarding_rule" "psc_ilb_consumer" {
  name   = "tf-test-psc-ilb-consumer-forwarding-rule%{random_suffix}"
  region = "us-central1"

  target                = google_compute_service_attachment.psc_ilb_service_attachment.id
  load_balancing_scheme = "" # need to override EXTERNAL default when target is a service attachment
  network               = "default"
  ip_address            = google_compute_address.psc_ilb_consumer_address.id
}

resource "google_compute_forwarding_rule" "psc_ilb_target_service" {
  name   = "tf-test-producer-forwarding-rule%{random_suffix}"
  region = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.producer_service_backend.id
  all_ports             = true
  network               = google_compute_network.psc_ilb_network.name
  subnetwork            = google_compute_subnetwork.psc_ilb_producer_subnetwork.name
}

resource "google_compute_region_backend_service" "producer_service_backend" {
  name   = "tf-test-producer-service%{random_suffix}"
  region = "us-central1"

  health_checks = [google_compute_health_check.producer_service_health_check.id]
}

resource "google_compute_health_check" "producer_service_health_check" {
  name = "tf-test-producer-service-health-check%{random_suffix}"

  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "80"
  }
}

resource "google_compute_network" "psc_ilb_network" {
  name = "tf-test-psc-ilb-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "psc_ilb_producer_subnetwork" {
  name   = "tf-test-psc-ilb-producer-subnetwork%{random_suffix}"
  region = "us-central1"

  network       = google_compute_network.psc_ilb_network.id
  ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_subnetwork" "psc_ilb_nat" {
  name   = "tf-test-psc-ilb-nat%{random_suffix}"
  region = "us-central1"

  network       = google_compute_network.psc_ilb_network.id
  purpose       =  "PRIVATE_SERVICE_CONNECT"
  ip_cidr_range = "10.1.0.0/16"
}

resource "google_compute_network" "consumer_net" {
  name                    = "tf-test-consumer-net%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "consumer_subnet" {
  name          = "tf-test-consumer-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.consumer_net.id
}

resource "google_compute_network" "consumer_net_2" {
  name                    = "tf-test-consumer-net-2%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "consumer_subnet_2" {
  name          = "tf-test-consumer-subnet-2%{random_suffix}"
  ip_cidr_range = "10.1.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.consumer_net.id
}
`, context)
}

func testAccNetworkConnectivityServiceConnectionMap_networkConnectivityMapBasicUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_service_connection_map" "default" {
  name = "tf-test-my-network-connectivity-map%{random_suffix}"
  description   = "my basic service connection map"
  location = "us-central1"
  service_class = "%{service_class_name}"
  producer_psc_configs {
    service_attachment_uri = google_compute_service_attachment.psc_ilb_service_attachment.id
  }
  consumer_psc_configs {
    network = google_compute_network.consumer_net_2.id
    project = "%{project_id}"
    consumer_instance_project = "%{project_id}"
  }
}

resource "google_compute_service_attachment" "psc_ilb_service_attachment" {
  name        = "tf-test-my-psc-ilb%{random_suffix}"
  region      = "us-central1"
  description = "A service attachment configured with Terraform"

  enable_proxy_protocol    = true
  connection_preference    = "ACCEPT_AUTOMATIC"
  nat_subnets              = [google_compute_subnetwork.psc_ilb_nat.id]
  target_service           = google_compute_forwarding_rule.psc_ilb_target_service.id
}

resource "google_compute_address" "psc_ilb_consumer_address" {
  name   = "tf-test-psc-ilb-consumer-address%{random_suffix}"
  region = "us-central1"

  subnetwork   = "default"
  address_type = "INTERNAL"
}

resource "google_compute_forwarding_rule" "psc_ilb_consumer" {
  name   = "tf-test-psc-ilb-consumer-forwarding-rule%{random_suffix}"
  region = "us-central1"

  target                = google_compute_service_attachment.psc_ilb_service_attachment.id
  load_balancing_scheme = "" # need to override EXTERNAL default when target is a service attachment
  network               = "default"
  ip_address            = google_compute_address.psc_ilb_consumer_address.id
}

resource "google_compute_forwarding_rule" "psc_ilb_target_service" {
  name   = "tf-test-producer-forwarding-rule%{random_suffix}"
  region = "us-central1"

  load_balancing_scheme = "INTERNAL"
  backend_service       = google_compute_region_backend_service.producer_service_backend.id
  all_ports             = true
  network               = google_compute_network.psc_ilb_network.name
  subnetwork            = google_compute_subnetwork.psc_ilb_producer_subnetwork.name
}

resource "google_compute_region_backend_service" "producer_service_backend" {
  name   = "tf-test-producer-service%{random_suffix}"
  region = "us-central1"

  health_checks = [google_compute_health_check.producer_service_health_check.id]
}

resource "google_compute_health_check" "producer_service_health_check" {
  name = "tf-test-producer-service-health-check%{random_suffix}"

  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "80"
  }
}

resource "google_compute_network" "psc_ilb_network" {
  name = "tf-test-psc-ilb-network%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "psc_ilb_producer_subnetwork" {
  name   = "tf-test-psc-ilb-producer-subnetwork%{random_suffix}"
  region = "us-central1"

  network       = google_compute_network.psc_ilb_network.id
  ip_cidr_range = "10.0.0.0/16"
}

resource "google_compute_subnetwork" "psc_ilb_nat" {
  name   = "tf-test-psc-ilb-nat%{random_suffix}"
  region = "us-central1"

  network       = google_compute_network.psc_ilb_network.id
  purpose       =  "PRIVATE_SERVICE_CONNECT"
  ip_cidr_range = "10.1.0.0/16"
}

resource "google_compute_network" "consumer_net" {
  name                    = "tf-test-consumer-net%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "consumer_subnet" {
  name          = "tf-test-consumer-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.consumer_net.id
}

resource "google_compute_network" "consumer_net_2" {
  name                    = "tf-test-consumer-net-2%{random_suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "consumer_subnet_2" {
  name          = "tf-test-consumer-subnet-2%{random_suffix}"
  ip_cidr_range = "10.1.0.0/24"
  region        = "us-central1"
  network       = google_compute_network.consumer_net.id
}
`, context)
}
