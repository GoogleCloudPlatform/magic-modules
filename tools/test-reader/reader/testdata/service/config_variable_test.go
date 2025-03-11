package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testConfigVariable = `
resource "config_variable" "basic" {
  field_one = "value-one"
}
`

func TestAccSqlDatabaseInstance_basicInferredName(t *testing.T) {
	VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testConfigVariable,
			},
		},
	})
}
