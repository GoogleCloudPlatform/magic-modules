package service_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

// A base config shared by several config funcs.
var testFormattingCallBase = `
resource "formatting_call_base" "base" {
  field_one = "value-one"
}
`

var testFormattingCallInline = `
resource "formatting_call_inline" "inline" {
  field_two = "%s"
}
`

func TestAccFormattingCall(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				// A formatting call written inline in the step rather than
				// wrapped in a config func.
				Config: fmt.Sprintf(testFormattingCallInline, "value-two"),
			},
			{
				// A config func whose template is a shared base config
				// concatenated onto a literal.
				Config: testAccFormattingCallConcat("value-three"),
			},
			{
				// A config func returning a call to another config func whose
				// first argument is a value to interpolate, not a config.
				Config: testAccFormattingCallNested(),
			},
			{
				// A config assembled from two config funcs in the step itself.
				Config: testAccFormattingCallStringArg("value-five") + testFormattingCallBase,
			},
			{
				// A config written as a literal in the step itself.
				Config: `
resource "formatting_call_literal" "literal" {
  field_six = "value-six"
}
`,
			},
		},
	})
}

func testAccFormattingCallConcat(fieldThree string) string {
	return fmt.Sprintf(testFormattingCallBase+`
resource "formatting_call_concat" "concat" {
  field_three = "%s"
}
`, fieldThree)
}

func testAccFormattingCallNested() string {
	return testAccFormattingCallStringArg("value-four")
}

func testAccFormattingCallStringArg(fieldFour string) string {
	return fmt.Sprintf(`
resource "formatting_call_string_arg" "string_arg" {
  field_four = "%s"
}
`, fieldFour)
}
