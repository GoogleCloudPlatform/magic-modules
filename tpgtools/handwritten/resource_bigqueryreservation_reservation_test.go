package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigqueryReservationReservation_bigqueryReservation(t *testing.T) {
	t.Parallel()

	location := "asia-northeast1"
	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
		"location":      location,
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryReservationReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservationReservation_bigqueryReservationBasic(context),
			},
			{
				ResourceName:      "google_bigquery_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigqueryReservationReservation_bigqueryReservationBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_bigquery_reservation" "reservation" {
	name           = "reservation%{random_suffix}"
	location       = "%{location}"
	// Set to 0 for testing purposes
	// In reality this would be larger than zero
	slot_capacity  = 0
	ignore_idle_slots = false
}
`, context)
}
