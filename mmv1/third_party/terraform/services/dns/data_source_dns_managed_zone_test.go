package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dns"
)

func TestAccDataSourceDnsManagedZone_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsManagedZone_basic(acctest.RandString(t, 10)),
				Check: acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
					"data.google_dns_managed_zone.qa",
					"google_dns_managed_zone.foo",
					[]string{
						"force_destroy",
					},
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZone_basic(managedZoneName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foo" {
  name        = "tf-test-qa-zone-%s"
  dns_name    = "qa.gcp.tfacc.hashicorptest.com."
  description = "QA DNS zone"
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.foo.name
}
`, managedZoneName)
}

func TestAccDataSourceDnsManagedZone_private(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsManagedZone_private(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_dns_managed_zone.qa",
						"google_dns_managed_zone.private-zone",
						[]string{
							"force_destroy",
						},
					),
					resource.TestCheckResourceAttr("data.google_dns_managed_zone.qa", "visibility", "private"),
					resource.TestCheckResourceAttr("data.google_dns_managed_zone.qa", "private_visibility_config.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_managed_zone.qa", "private_visibility_config.0.networks.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZone_private(suffix string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network-1" {
  name                    = "tf-test-network-1-%s"
  auto_create_subnetworks = false
}

resource "google_dns_managed_zone" "private-zone" {
  name        = "tf-test-private-zone-%s"
  dns_name    = "private.example.com."
  description = "Example private DNS zone"
  visibility  = "private"

  private_visibility_config {
    networks {
      network_url = google_compute_network.network-1.id
    }
  }
}

data "google_dns_managed_zone" "qa" {
  name = google_dns_managed_zone.private-zone.name
}
`, suffix, suffix)
}
