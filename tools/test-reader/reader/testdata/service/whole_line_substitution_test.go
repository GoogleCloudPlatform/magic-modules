package service_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccWholeLineSubstitution(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccWholeLineSubstitution("value-two"),
			},
		},
	})
}

// Some configs use a substitution that occupies an entire line to inject whole
// HCL statements rather than a value: a shared setup block at the top level, or
// an attribute or nested block that is only present in some test cases. Those
// lines are dropped, so nothing they inject counts as covered, but every other
// field in the config is still read.
func testAccWholeLineSubstitution(fieldTwo string) string {
	return fmt.Sprintf(`
%s

resource "whole_line_substitution" "resource" {
  field_one = "value-one"
  %s
  field_three {
    %s
    field_five = "value-five"
  }
}
`,
		testAccWholeLineSubstitutionSetup(),
		fmt.Sprintf("field_two = %q", fieldTwo),
		`field_four = "value-four"`,
	)
}

func testAccWholeLineSubstitutionSetup() string {
	return `
resource "whole_line_substitution_setup" "setup" {
  field_one = "value-one"
}
`
}
