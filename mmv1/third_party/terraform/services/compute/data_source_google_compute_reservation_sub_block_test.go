package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeReservationSubBlock_basic(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	dsName := "my_reservation_sub_block"
	dsFullName := fmt.Sprintf("data.google_compute_reservation_sub_block.%s", dsName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeReservationSubBlockConfig(reservationName, dsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsFullName, "name"),
					resource.TestCheckResourceAttrSet(dsFullName, "zone"),
					resource.TestCheckResourceAttrSet(dsFullName, "project"),
					resource.TestCheckResourceAttrSet(dsFullName, "kind"),
					resource.TestCheckResourceAttrSet(dsFullName, "self_link"),
					resource.TestCheckResourceAttrSet(dsFullName, "sub_block_count"),
					resource.TestCheckResourceAttrSet(dsFullName, "status"),
				),
			},
		},
	})
}

func testAccDataSourceComputeReservationSubBlockConfig(reservationName, dsName string) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-west1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}

# Note: Reservation sub-blocks are automatically created by Google Cloud
# This data source would reference an existing sub-block under a reservation block
# In a real scenario, you would need to query the API to get the actual block and sub-block names
data "google_compute_reservation_sub_block" "%s" {
  name              = "sub-block-name-from-api"
  reservation_block = "block-name-from-api"
  reservation       = google_compute_reservation.reservation.name
  zone              = "us-west1-a"
  
  depends_on = [google_compute_reservation.reservation]
}
`, reservationName, dsName)
}
