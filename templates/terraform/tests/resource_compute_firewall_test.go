package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeFirewall_basic(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_update(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeFirewall_update(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_priority(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFirewall_priority(networkName, firewallName, 1001),
			},
			{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_noSource(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_noSource(networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_denied(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_denied(networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_egress(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_egress(networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_serviceAccounts(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	sourceSa := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	targetSa := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_serviceAccounts(sourceSa, targetSa, networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeFirewall_disabled(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))
	firewallName := fmt.Sprintf("firewall-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeFirewall_disabled(networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccComputeFirewall_basic(networkName, firewallName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_firewall.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckComputeFirewallDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_firewall" {
			continue
		}

		_, err := config.clientCompute.Firewalls.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("Firewall still exists")
		}
	}

	return nil
}

func testAccComputeFirewall_basic(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		allow {
			protocol = "icmp"
		}
	}`, network, firewall)
}

func testAccComputeFirewall_update(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "New description"
		network = "${google_compute_network.foobar.self_link}"
		source_service_accounts = ["sa@fake-project.iam.gserviceaccount.com"]
		priority = 2000

		deny {
			protocol = "tcp"
			ports = ["80-255"]
		}

		disabled = true
	}`, network, firewall)
}

func testAccComputeFirewall_priority(network, firewall string, priority int) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		allow {
			protocol = "icmp"
		}
		priority = %d
	}`, network, firewall, priority)
}

func testAccComputeFirewall_noSource(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"

		allow {
			protocol = "tcp"
			ports    = [22]
		}
	}`, network, firewall)
}

func testAccComputeFirewall_denied(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		deny {
			protocol = "tcp"
			ports    = [22]
		}
	}`, network, firewall)
}

func testAccComputeFirewall_egress(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		direction = "EGRESS"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"

		allow {
			protocol = "tcp"
			ports    = [22]
		}
	}`, network, firewall)
}

func testAccComputeFirewall_serviceAccounts(sourceSa, targetSa, network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_service_account" "source" {
		account_id = "%s"
	}

	resource "google_service_account" "target" {
		account_id = "%s"
	}

	resource "google_compute_network" "foobar" {
		name = "%s"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"

		allow {
			protocol = "icmp"
		}

		source_service_accounts = ["${google_service_account.source.email}"]
		target_service_accounts = ["${google_service_account.target.email}"]
	}`, sourceSa, targetSa, network, firewall)
}

func testAccComputeFirewall_disabled(network, firewall string) string {
	return fmt.Sprintf(`
	resource "google_compute_network" "foobar" {
		name = "%s"
		auto_create_subnetworks = false
		ipv4_range = "10.0.0.0/16"
	}

	resource "google_compute_firewall" "foobar" {
		name = "%s"
		description = "Resource created for Terraform acceptance testing"
		network = "${google_compute_network.foobar.name}"
		source_tags = ["foo"]

		allow {
			protocol = "icmp"
		}

		disabled = true
	}`, network, firewall)
}
