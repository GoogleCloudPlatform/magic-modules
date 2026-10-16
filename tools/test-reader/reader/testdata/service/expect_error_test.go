package service_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccExpectError(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccExpectError_basic(),
			},
			{
				Config:      testAccExpectError_invalid(),
				ExpectError: regexp.MustCompile("invalid"),
			},
			{
				Config:   testAccExpectError_planOnly(),
				PlanOnly: true,
			},
		},
	})
}

func testAccExpectError_basic() string {
	return `
resource "expect_error_resource" "resource" {
  field_one = "value-one"
}
`
}

func testAccExpectError_invalid() string {
	return `
resource "expect_error_resource" "resource" {
  field_one = "value-one"
  field_two = "invalid-value"
}
`
}

func testAccExpectError_planOnly() string {
	return `
resource "expect_error_resource" "resource" {
  field_one   = "value-one"
  field_three = "plan-only-value"
}
`
}
