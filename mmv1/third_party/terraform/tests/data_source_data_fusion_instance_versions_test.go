package google

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDataFusionInstanceVersions_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataFusionInstanceVersions_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDataFusionInstanceVersionsMeta("data.google_data_fusion_instance_versions.versions"),
				),
			},
		},
	})
}

func testAccCheckGoogleDataFusionInstanceVersionsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find versions data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("versions data source ID not set.")
		}

		versionCountStr, ok := rs.Primary.Attributes["instance_versions.#"]
		if !ok {
			return errors.New("can't find 'image_versions' attribute")
		}

		versionCount, err := strconv.Atoi(versionCountStr)
		if err != nil {
			return errors.New("failed to read number of valid instance versions")
		}
		if versionCount < 1 {
			return fmt.Errorf("expected at least 1 valid instance versions, received %d, this is most likely a bug",
				versionCount)
		}

		for i := 0; i < versionCount; i++ {
			idx := "instance_versions." + strconv.Itoa(i)
			if v, ok := rs.Primary.Attributes[idx+".version_number"]; !ok || v == "" {
				return fmt.Errorf("instance_version %v is missing version_number", i)
			}
		}

		return nil
	}
}

var testAccDataSourceDataFusionInstanceVersions_basic = `
data "google_data_fusion_instance_versions" "versions" {
  location = "us-central1"
}
`
