package google

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAlloydbLocation_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlloydbLocation_basic(context),
				Check: resource.ComposeTestCheckFunc(
					validateAlloydbLocationResult(
						"data.google_alloydb_location.qa",
					),
				),
			},
		},
	})
}

func testAccDataSourceAlloydbLocation_basic(context map[string]interface{}) string {
	return Nprintf(`
data "google_alloydb_location" "qa"{
	location = "us-central1"
}
`, context)
}

func validateAlloydbLocationResult(dataSourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}
		var dsAttr map[string]string
		dsAttr = ds.Primary.Attributes
		if dsAttr["name"] == "" {
			return errors.New("name parameter is not set for the location")
		}
		if dsAttr["location_id"] == "" {
			return errors.New("location_id parameter is not set for the location")
		}
		return nil
	}
}
