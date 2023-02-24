package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMultipleResources(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccMultipleResources(),
			},
			{
				ImportStateVerify: true,
			},
			{
				Config: testAccMultipleResources_update(),
			},
		},
	})
}

func testAccMultipleResources() string {
	return Nprintf(`
resource "resource_one" "instance_one" {
  field_one = "value-one"
}

resource "resource_one" "instace_two" {
  field_one = "value-one"
}

resource "resource_two" "instace_one" {
  field_one = "value-one"
}

resource "resource_two" "instace_two" {
  field_one = "value-one"
}`, context)
}

func testAccMultipleResources_update() string {
	return Nprintf(`
resource "resource_one" "instance_one" {
  field_one = "value-two"
}

resource "resource_one" "instace_two" {
  field_one = "value-two"
}

resource "resource_two" "instace_one" {
  field_one = "value-two"
}

resource "resource_two" "instace_two" {
  field_one = "value-two"
}`, context)
}
