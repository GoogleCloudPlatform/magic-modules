package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCoveredResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccCoveredResource(),
			},
			{
				Config: testAccCoveredResource_update(),
			},
		},
	})
}

func testAccCoveredResource() string {
	return Nprintf(`
resource "covered_resource" "resource" {
  field_one = "value-one"
  field_four {
    field_five {
      field_six = "value-three"
    }
  }
}
`, context)
}

func testAccCoveredResource_update() string {
	return Nprintf(`
resource "covered_resource" "resource" {
  field_two {
    field_three = "value-two"
  }
  field_four {
    field_five {
      field_six = "value-three"
    }
  }
}
`, context)
}
