package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccDataSourceGoogleNetwork(t *testing.T) {
	t.Parallel()

	networkName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleNetworkConfig(networkName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleNetworkCheck("data.google_compute_network.my_network", "google_compute_network.foobar"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleNetwork_peerings(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	dsName := "data.google_compute_network.network_a"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleNetwork_peeringsConfig(suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "peerings.#", "1"),
					resource.TestCheckResourceAttrPair(dsName, "peerings.0.name", "google_compute_network_peering.a_to_b", "name"),
					resource.TestCheckResourceAttrPair(dsName, "peerings.0.network", "google_compute_network.network_b", "self_link"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(dsName, "peerings.0.state_details"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.export_custom_routes", "true"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.import_custom_routes", "false"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.export_subnet_routes_with_public_ip", "true"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.import_subnet_routes_with_public_ip", "false"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.stack_type", "IPV4_ONLY"),
					resource.TestCheckResourceAttr(dsName, "peerings.0.update_strategy", "INDEPENDENT"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleNetworkCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes
		network_attrs_to_test := []string{
			"id",
			"name",
			"network_id",
			"numeric_id",
			"description",
			"internal_ipv6_range",
		}

		for _, attr_to_check := range network_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if !tpgresource.CompareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		return nil
	}
}

func testAccDataSourceGoogleNetworkConfig(name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                     = "%s"
  description              = "my-description"
  enable_ula_internal_ipv6 = true
  auto_create_subnetworks  = false
}

data "google_compute_network" "my_network" {
  name = google_compute_network.foobar.name
}
`, name)
}

func testAccDataSourceGoogleNetwork_peeringsConfig(suffix string) string {
	return acctest.Nprintf(`
resource "google_compute_network" "network_a" {
  name                    = "tf-test-network-a-%{suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network" "network_b" {
  name                    = "tf-test-network-b-%{suffix}"
  auto_create_subnetworks = false
}

resource "google_compute_network_peering" "a_to_b" {
  name                 = "tf-test-peering-a-b-%{suffix}"
  network              = google_compute_network.network_a.self_link
  peer_network         = google_compute_network.network_b.self_link
  export_custom_routes = true
}

# The reverse peering makes the a_to_b peering ACTIVE.
resource "google_compute_network_peering" "b_to_a" {
  name         = "tf-test-peering-b-a-%{suffix}"
  network      = google_compute_network.network_b.self_link
  peer_network = google_compute_network.network_a.self_link

  depends_on = [google_compute_network_peering.a_to_b]
}

data "google_compute_network" "network_a" {
  name = google_compute_network.network_a.name

  depends_on = [
    google_compute_network_peering.a_to_b,
    google_compute_network_peering.b_to_a,
  ]
}
`, map[string]interface{}{"suffix": suffix})
}
