package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceComputeReservationHost_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"host_name":     fmt.Sprintf("tf-test-host-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceComputeReservationHost_topLevel(context),
				ExpectError: regexp.MustCompile(`Error reading Host .*: googleapi: Error 404`),
			},
			{
				Config:      testAccDataSourceComputeReservationHost_hierarchical(context),
				ExpectError: regexp.MustCompile(`Error reading Host .*: googleapi: Error (400|404)`),
			},
			{
				Config:      testAccDataSourceComputeReservationHost_missingReservation(context),
				ExpectError: regexp.MustCompile("(?s)all of.*`reservation,reservation_block`.*must be specified"),
			},
		},
	})
}

func testAccDataSourceComputeReservationHost_topLevel(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_host" "top" {
  name = "%{host_name}"
  zone = "us-central1-a"
}
`, context)
}

func testAccDataSourceComputeReservationHost_hierarchical(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_host" "scoped" {
  name                  = "%{host_name}"
  reservation           = "tf-test-res-%{random_suffix}"
  reservation_block     = "tf-test-blk-%{random_suffix}"
  reservation_sub_block = "tf-test-sub-%{random_suffix}"
  zone                  = "us-central1-a"
}
`, context)
}

func testAccDataSourceComputeReservationHost_missingReservation(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_host" "invalid" {
  name              = "%{host_name}"
  reservation_block = "tf-test-blk-%{random_suffix}"
  zone              = "us-central1-a"
}
`, context)
}
