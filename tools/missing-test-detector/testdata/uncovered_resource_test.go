package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
