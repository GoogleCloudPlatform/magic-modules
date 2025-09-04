package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleComputeInterconnectLocation_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeInterconnectLocation_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_compute_interconnect_location.iad_zone1", "self_link"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleComputeInterconnectLocation_basic() string {
	return `
data "google_compute_interconnect_location" "iad_zone1" {
	name = "iad-zone1-1"
}
`
}
