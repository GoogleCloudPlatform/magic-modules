package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// EnableModel is a beta-only, synchronous, action-only resource. It has no
// Read, Delete, or Import support, so this test only verifies that an apply
// succeeds and that the output-only fields are populated. There is no inverse
// "disable" API, so destroying the resource is a no-op and CheckDestroy is
// intentionally omitted.
func TestAccVertexAIModelGardenEnableModel_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIModelGardenEnableModel_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"google_vertex_ai_model_garden_enable_model.enable", "enablement_state"),
					resource.TestCheckResourceAttrSet(
						"google_vertex_ai_model_garden_enable_model.enable", "publisher_endpoint"),
				),
			},
		},
	})
}

func testAccVertexAIModelGardenEnableModel_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_model_garden_enable_model" "enable" {
  provider             = google-beta
  project              = data.google_project.project.project_id
  publisher_model_name = "publishers/google/models/paligemma@paligemma-224-float32"
}

data "google_project" "project" {
  provider = google-beta
}
`, context)
}
