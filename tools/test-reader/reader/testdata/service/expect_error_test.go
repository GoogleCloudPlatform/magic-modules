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
			{
				Config:  testAccExpectError_destroy(),
				Destroy: true,
			},
			{
				Config:      testAccExpectError_importState(),
				ImportState: true,
			},
			{
				Config:       testAccExpectError_refreshState(),
				RefreshState: true,
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

func testAccExpectError_destroy() string {
	return `
resource "expect_error_resource" "resource" {
  field_one  = "value-one"
  field_four = "destroy-value"
}
`
}

func testAccExpectError_importState() string {
	return `
resource "expect_error_resource" "resource" {
  field_one  = "value-one"
  field_five = "import-value"
}
`
}

func testAccExpectError_refreshState() string {
	return `
resource "expect_error_resource" "resource" {
  field_one = "value-one"
  field_six = "refresh-value"
}
`
}
