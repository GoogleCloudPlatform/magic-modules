package compute_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceComputeHostVersion_basic(t *testing.T) {
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
				Config:      testAccDataSourceComputeHostVersion_topLevel(context),
				ExpectError: regexp.MustCompile(`Error getting version for Host .*: googleapi: Error (400|404)`),
			},
			{
				Config:      testAccDataSourceComputeHostVersion_withSelections(context),
				ExpectError: regexp.MustCompile(`Error getting version for Host .*: googleapi: Error (400|404)`),
			},
			{
				Config:      testAccDataSourceComputeHostVersion_missingBlock(context),
				ExpectError: regexp.MustCompile("(?s)all of.*`reservation,reservation_block,reservation_sub_block`.*must be specified"),
			},
		},
	})
}

func testAccDataSourceComputeHostVersion_topLevel(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_host_version" "top" {
  name = "%{host_name}"
  zone = "us-central1-a"
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
  sbom_selections       = ["SBOM_SELECTION_CURRENT"]
}
`, context)
}

func testAccDataSourceComputeHostVersion_missingBlock(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_host_version" "invalid" {
  name                  = "%{host_name}"
  reservation           = "tf-test-res-%{random_suffix}"
  reservation_sub_block = "tf-test-sub-%{random_suffix}"
  zone                  = "us-central1-a"
}
`, context)
}
