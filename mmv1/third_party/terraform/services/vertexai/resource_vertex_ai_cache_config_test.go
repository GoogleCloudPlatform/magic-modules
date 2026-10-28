package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/vertexai"
)

func TestAccVertexAICacheConfig_region(t *testing.T) {
	// The cache config is a project-level singleton, so tests that modify it
	// must not run in parallel with each other.
	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"region":        "europe-west4",
		"disable_cache": false,
	}
	contextUpdated := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"region":        "asia-northeast1",
		"disable_cache": true,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAICacheConfig_region(context),
			},
			{
				ResourceName:            "google_vertex_ai_cache_config.cache_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
			{
				// Changing the region only changes the API endpoint, so it is an
				// in-place update rather than a replacement.
				Config: testAccVertexAICacheConfig_region(contextUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_vertex_ai_cache_config.cache_config", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_vertex_ai_cache_config.cache_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"region"},
			},
		},
	})
}

func testAccVertexAICacheConfig_region(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_cache_config" "cache_config" {
  project       = "%{project}"
  region        = "%{region}"
  disable_cache = %{disable_cache}
}
`, context)
}
