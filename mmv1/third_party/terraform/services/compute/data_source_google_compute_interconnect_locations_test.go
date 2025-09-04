package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleComputeInterconnectLocations_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeInterconnectLocations_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_compute_interconnect_locations.all", "locations.0.self_link"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeInterconnectLocations_basic() string {
	return `
data "google_compute_interconnect_locations" "all" {}
`
}
