package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccMultipleResources(t *testing.T) {
	VcrTest(t, resource.TestCase{
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
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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
