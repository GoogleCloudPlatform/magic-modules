package service_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccEscapedPercent(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccEscapedPercentSprintf("value-four"),
			},
			{
				Config: testAccEscapedPercentNprintf(map[string]interface{}{"suffix": "value"}),
			},
			{
				Config: testAccEscapedPercentNprintfSprintf(map[string]interface{}{"suffix": "value"}, "value-two"),
			},
		},
	})
}

// fmt.Sprintf reads "%%" as an escaped percent sign.
func testAccEscapedPercentSprintf(fieldFour string) string {
	return fmt.Sprintf(`
resource "escaped_percent" "resource" {
  field_one   = parseint("ff", 16) %% pow(2, 32)
  field_two   = "100%%"
  field_three = "%%d"
  field_four  = "%s"
}
`, fieldFour)
}

// acctest.Nprintf only replaces %{key}, so "%%" is two percent signs.
func testAccEscapedPercentNprintf(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "escaped_percent" "resource" {
  field_one = "name-%%{suffix}"
}
`, context)
}

// fmt.Sprintf formats "%%{suffix}" to "%{suffix}", which acctest.Nprintf then
// replaces.
func testAccEscapedPercentNprintfSprintf(context map[string]interface{}, fieldTwo string) string {
	return acctest.Nprintf(fmt.Sprintf(`
resource "escaped_percent" "resource" {
  field_one = "name-%%{suffix}"
  field_two = "%s"
}
`, fieldTwo), context)
}
