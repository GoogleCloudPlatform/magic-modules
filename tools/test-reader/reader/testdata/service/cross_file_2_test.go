package service_test

import (
	"google/provider/new/google-beta/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCrossFile2(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccCrossFileConfig2(),
			},
		},
	})
}

func testAccCrossFileConfig1() string {
	return acctest.Nprintf(`
resource "serial_resource" "resource" {
  field_one = "value-one"
}
`, context)
}

func testAccCrossFileConfig2() string {
	return acctest.Nprintf(`
resource "serial_resource" "resource" {
  field_two {
    field_three = "value-two"
  }
}
`, context)
}
