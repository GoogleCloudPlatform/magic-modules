package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVertexAIRagEngineConfig_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project": envvar.GetTestProjectFromEnv(),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIRagEngineConfig_basic(context),
			},
		},
	})
}

func testAccVertexAIRagEngineConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_rag_engine_config" "test" {
  region = "us-central1"
  rag_managed_db_config {
    basic {}
  }
}
`, context)
}
