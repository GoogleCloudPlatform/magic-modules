package service_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccCoveredResource(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
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
	return fmt.Sprintf(acctest.Nprintf(`
resource "covered_resource" "resource" {
  field_one = "value-one"
  field_four {
    field_five {
      field_six = %v
    }
  }`), 0) + acctest.Nprintf(`
  field_seven = %{bool}
}
`, context)
}

func testAccCoveredResource_update() string {
	return `
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
`
}
