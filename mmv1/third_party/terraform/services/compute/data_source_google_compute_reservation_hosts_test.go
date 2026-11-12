package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceComputeReservationHosts_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"reservation_name": fmt.Sprintf("tf-test-res-hosts-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceComputeReservationHosts_missingReservation(context),
				ExpectError: regexp.MustCompile("(?s)all of.*`reservation,reservation_block`.*must be specified"),
			},
			{
				Config: testAccDataSourceComputeReservationHosts_topLevel(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_compute_reservation_hosts.top", "id", regexp.MustCompile(`^projects/[^/]+/zones/us-central1-a/hosts$`)),
					resource.TestCheckResourceAttr("data.google_compute_reservation_hosts.top", "zone", "us-central1-a"),
					resource.TestCheckResourceAttrSet("data.google_compute_reservation_hosts.top", "project"),
					resource.TestCheckResourceAttrSet("data.google_compute_reservation_hosts.top", "hosts.#"),
				),
			},
			{
				Config: testAccDataSourceComputeReservationHosts_withReservation(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_compute_reservation_hosts.scoped", "id", regexp.MustCompile(`^projects/[^/]+/zones/us-central1-c/reservations/tf-test-res-hosts-[^/]+/hosts$`)),
					resource.TestCheckResourceAttr("data.google_compute_reservation_hosts.scoped", "zone", "us-central1-c"),
					resource.TestCheckResourceAttr("data.google_compute_reservation_hosts.scoped", "hosts.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceComputeReservationHosts_topLevel(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_hosts" "top" {
  zone = "us-central1-a"
}
`, context)
}

func testAccDataSourceComputeReservationHosts_withReservation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%{reservation_name}"
  zone = "us-central1-c"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type = "n2-standard-2"
    }
  }
}

data "google_compute_reservation_hosts" "scoped" {
  reservation = google_compute_reservation.reservation.name
  zone        = "us-central1-c"
  filter      = "state = ACTIVE"
}
`, context)
}

func testAccDataSourceComputeReservationHosts_missingReservation(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_hosts" "invalid" {
  reservation_block     = "my-block"
  reservation_sub_block = "my-sub-block"
  zone                  = "us-central1-a"
}
`, context)
}
