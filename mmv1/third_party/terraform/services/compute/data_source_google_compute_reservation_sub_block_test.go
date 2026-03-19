package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeReservationSubBlock_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeReservationSubBlockConfig(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_reservation_sub_block.default", "reservation", fmt.Sprintf("tf-test-res-%s", context["random_suffix"])),
					resource.TestCheckResourceAttrSet("data.google_compute_reservation_sub_block.default", "reservation_block"),
					resource.TestCheckResourceAttrSet("data.google_compute_reservation_sub_block.default", "name"),
					resource.TestCheckResourceAttr("data.google_compute_reservation_sub_block.default", "zone", "us-central1-a"),
					resource.TestCheckResourceAttr("data.google_compute_reservation_sub_block.default", "sub_block_count", "1"),
					resource.TestCheckResourceAttr("data.google_compute_reservation_sub_block.default", "in_use_count", "0"),
				),
			},
		},
	})
}

func testAccDataSourceComputeReservationSubBlockConfig(context map[string]interface{}) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "tf-test-res-%s"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      machine_type = "a3-highgpu-8g"
      guest_accelerators {
        accelerator_type  = "nvidia-h100-80gb"
        accelerator_count = 8
      }
    }
  }
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_compute_reservation.reservation]

  create_duration = "120s"
}

data "google_compute_reservation_block" "default" {
  name        = google_compute_reservation.reservation.block_names[0]
  reservation = google_compute_reservation.reservation.name
  zone        = "us-central1-a"

  depends_on = [time_sleep.wait_120_seconds]
}

data "google_compute_reservation_sub_block" "default" {
  name              = data.google_compute_reservation_block.default.sub_block_names[0]
  reservation_block = data.google_compute_reservation_block.default.name
  reservation       = google_compute_reservation.reservation.name
  zone              = "us-central1-a"

  depends_on = [time_sleep.wait_120_seconds]
}
`, context["random_suffix"])
}
