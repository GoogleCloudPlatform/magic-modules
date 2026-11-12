package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceComputeReservationSubBlockVersion_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceComputeReservationSubBlockVersion_omittedSelections(context),
				ExpectError: regexp.MustCompile(`Error getting version for ReservationSubBlock .*: googleapi: Error (400|404)`),
			},
			{
				Config:      testAccDataSourceComputeReservationSubBlockVersion_withSelections(context),
				ExpectError: regexp.MustCompile(`Error getting version for ReservationSubBlock .*: googleapi: Error (400|404)`),
			},
		},
	})
}

func testAccDataSourceComputeReservationSubBlockVersion_omittedSelections(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_sub_block_version" "default" {
  name              = "tf-test-sub-%{random_suffix}"
  reservation_block = "tf-test-blk-%{random_suffix}"
  reservation       = "tf-test-res-%{random_suffix}"
  zone              = "us-central1-a"
}
`, context)
}

func testAccDataSourceComputeReservationSubBlockVersion_withSelections(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_reservation_sub_block_version" "selected" {
  name              = "tf-test-sub-%{random_suffix}"
  reservation_block = "tf-test-blk-%{random_suffix}"
  reservation       = "tf-test-res-%{random_suffix}"
  zone              = "us-central1-a"
  sbom_selections   = ["SBOM_SELECTION_CURRENT"]
}
`, context)
}
