package lustre_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccLustreInstance_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLustreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLustreInstance_full(context),
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccLustreInstance_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_lustre_instance.description", plancheck.ResourceActionUpdate),
					}
				}
			},
			{
				ResourceName:            "google_lustre_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccLustreInstance_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id = "instance%{random_suffix}"
  location = "us-central1-a"
  description = "test instance"
  filesystem = "testfs
  capacity_gib = 14000
  network = google_compute_network.network.name
  labels = {
    test = "value"
  }
  depends_on = [google_service_networking_connection.default]
}

resource "google_compute_network" "network" {
  name                    = "network%{random_suffix}"
  auto_create_subnetworks = true
  mtu = 8896
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name          = ""
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = google_compute_network.network.id
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

func testAccLustreInstance_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_lustre_instance" "instance" {
  instance_id = "instance%{random_suffix}"
  location = "us-central1-a"
  description = "test instance description field has been updated"
  filesystem = "testfs
  capacity_gib = 14000
  network = google_compute_network.network.name
  labels = {
    test = "value"
  }
  depends_on = [google_service_networking_connection.default]
}

resource "google_compute_network" "network" {
  name                    = "network%{random_suffix}"
  auto_create_subnetworks = true
  mtu = 8896
}

# Create an IP address
resource "google_compute_global_address" "private_ip_alloc" {
  name          = ""
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 24
  network       = google_compute_network.network.id
}

# Create a private connection
resource "google_service_networking_connection" "default" {
  network                 = google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_alloc.name]
}
`, context)
}

func testAccCheckLustreInstanceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_lustre_instance" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{LustreBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance_id}}")
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
				return fmt.Errorf("LustreInstance still exists at %s", url)
			}
		}

		return nil
	}
}
