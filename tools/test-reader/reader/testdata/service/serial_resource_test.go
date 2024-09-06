package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccSerialResource(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"test1": testAccSerialResource1,
		"test2": testAccSerialResource2,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccSerialResource1(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccSerialResourceConfig1(),
			},
		},
	})
}

func testAccSerialResource2(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccSerialResourceConfig2(),
			},
		},
	})
}

func testAccSerialResourceConfig1() string {
	return acctest.Nprintf(`
resource "serial_resource" "resource" {
  field_one = "value-one"
}
`, context)
}

func testAccSerialResourceConfig2() string {
	return acctest.Nprintf(`
resource "serial_resource" "resource" {
  field_two {
    field_three = "value-two"
  }
}
`, context)
}
