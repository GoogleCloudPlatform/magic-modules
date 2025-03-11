package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccFunctionCallResource(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionCallResource(),
			},
		},
	})
}

func helperFunction() string {
	return `
resource "helper_resource" "default" {
  field_one = "value-one"
}
`
}

func testAccFunctionCallResource() string {
	return helperFunction() + acctest.Nprintf(`
resource "helped_resource" "primary" {
  field_one = "value-one"
}
`)
}
