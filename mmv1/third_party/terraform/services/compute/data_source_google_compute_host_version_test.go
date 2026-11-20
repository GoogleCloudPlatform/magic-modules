package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

// Compute Engine has no consumer hosts.insert API; compute#host entities with SBOM
// version reporting are system-populated only on physical Extended Reservations
// provisioned with dedicated host/SBOM hardware entitlements that cannot be
// dynamically allocated in standard CI projects. These acceptance steps verify
// live REST endpoint wiring and request payload validation against /compute/v1/
// and /compute/beta/ (returning 404 for a non-existent host), while unit tests in
// data_source_google_compute_host_version_internal_test.go cover 200/LRO polling,
// inlineSbomInfo flattening, and client-side schema validation.
func TestAccDataSourceComputeHostVersion_basic(t *testing.T) {
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
				Config:      testAccDataSourceComputeHostVersion_topLevel(context),
				ExpectError: regexp.MustCompile(`Error getting version for Host .*: googleapi: Error 404.*notFound`),
			},
			{
				Config:      testAccDataSourceComputeHostVersion_withSelections(context),
				ExpectError: regexp.MustCompile(`Error getting version for Host .*: googleapi: Error 404.*notFound`),
			},
		},
	})
}

func testAccDataSourceComputeHostVersion_topLevel(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_host_version" "top" {
  name            = "%{host_name}"
  zone            = "us-central1-a"
  sbom_selections = ["SBOM_SELECTION_CURRENT"]
}
`, context)
}

func testAccDataSourceComputeHostVersion_withSelections(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_host_version" "scoped" {
  name                  = "%{host_name}"
  reservation           = "tf-test-res-%{random_suffix}"
  reservation_block     = "tf-test-blk-%{random_suffix}"
  reservation_sub_block = "tf-test-sub-%{random_suffix}"
  zone                  = "us-central1-a"
  sbom_selections       = ["SBOM_SELECTION_TARGET"]
}
`, context)
}
