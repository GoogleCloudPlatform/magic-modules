package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceGoogleComputeRoutes_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-routes-%d", acctest.RandInt(t))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeRouteDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeRoutesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.#", "1"),
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.0.name", name),
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.0.dest_range", "15.0.0.0/24"),
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.0.priority", "100"),
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.0.tags.#", "1"),
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.0.tags.0", "tf-test-tag"),
					resource.TestCheckResourceAttr("data.google_compute_routes.one", "routes.0.route_type", "STATIC"),
					resource.TestCheckResourceAttrSet("data.google_compute_routes.one", "routes.0.creation_timestamp"),
					resource.TestCheckResourceAttrPair("data.google_compute_routes.one", "routes.0.network", "google_compute_route.foobar", "network"),
					// next_hop_gateway is only checked as set, because google_compute_route keeps the API's form of that link.
					resource.TestCheckResourceAttrSet("data.google_compute_routes.one", "routes.0.next_hop_gateway"),
					resource.TestCheckResourceAttrPair("data.google_compute_routes.one", "routes.0.self_link", "google_compute_route.foobar", "self_link"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeRoutesConfig(name string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "foobar" {
  name                    = "%s"
  auto_create_subnetworks = false
}

resource "google_compute_route" "foobar" {
  name             = "%s"
  dest_range       = "15.0.0.0/24"
  network          = google_compute_network.foobar.name
  next_hop_gateway = "default-internet-gateway"
  priority         = 100
  tags             = ["tf-test-tag"]
}

data "google_compute_routes" "one" {
  filter = "name eq ${google_compute_route.foobar.name}"
}
`, name, name)
}
