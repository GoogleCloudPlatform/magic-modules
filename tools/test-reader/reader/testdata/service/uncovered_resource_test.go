package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccUncoveredResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccUncoveredResource(),
			},
		},
	})
}

func testAccUncoveredResource() string {
	return acctest.Nprintf(`
resource "uncovered_resource" "resource" {
  field_two {
    field_three = "value-two"
  }
}
`, context)
}
