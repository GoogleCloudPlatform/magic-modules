package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testConfigVariable = `
resource "config_variable" "basic" {
  field_one = "value-one"
}
`

func TestAccSqlDatabaseInstance_basicInferredName(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testConfigVariable,
			},
		},
	})
}
