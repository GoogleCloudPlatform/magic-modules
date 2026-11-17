package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

// Compute Engine has no consumer hosts.insert API; compute#host entities are
// system-populated only on physical Extended Reservations provisioned with
// dedicated host hardware entitlements that cannot be dynamically allocated in
// standard CI projects. These acceptance steps verify live REST endpoint
// wiring against /compute/v1/ and /compute/beta/ (returning 404 for a
// non-existent host), while unit tests in
// data_source_google_compute_reservation_host_internal_test.go cover 200 OK
// response flattening, canonical name write-back, and schema validation.
func TestAccDataSourceComputeReservationHost_basic(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	context := map[string]interface{}{
		"random_suffix": suffix,
		"host_name":     fmt.Sprintf("tf-test-host-%s", suffix),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceComputeReservationHost_topLevel(context),
				ExpectError: regexp.MustCompile(`Error reading Host .*: googleapi: Error 404.*notFound`),
			},
			{
				Config:      testAccDataSourceComputeReservationHost_hierarchical(context),
				ExpectError: regexp.MustCompile(`Error reading Host .*: googleapi: Error 404.*notFound`),
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
