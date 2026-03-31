package dialogflow_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccDialogflowTool_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"tool_name":     "tf-test-tool-" + acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDialogflowToolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDialogflowTool_basic(context),
			},
			{
				ResourceName:            "google_dialogflow_tool.basic_tool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccDialogflowTool_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dialogflow_tool.basic_tool", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dialogflow_tool.basic_tool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccDialogflowTool_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_tool" "basic_tool" {
  location = "global"
  display_name = "%{tool_name}"
  description = "A basic open_api_spec tool"
  tool_key = "%{tool_name}"
  open_api_spec {
    text_schema = "openapi: 3.0.0\ninfo:\n  title: Example API\n  version: 1.0.0\npaths:\n  /example:\n    get:\n      summary: Example GET\n      responses:\n        '200':\n          description: OK"
  }
}
`, context)
}

func testAccDialogflowTool_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dialogflow_tool" "basic_tool" {
  location = "global"
  display_name = "%{tool_name}-updated"
  description = "A basic open_api_spec tool updated"
	tool_key = "%{tool_name}-updated"
  open_api_spec {
    text_schema = "openapi: 3.0.0\ninfo:\n  title: Example API Updated\n  version: 1.0.0\npaths:\n  /example:\n    get:\n      summary: Example GET Updated\n      responses:\n        '200':\n          description: OK"
  }
}
`, context)
}
