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
		"project": envvar.GetTestProjectFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAICacheConfig_region(context, "europe-west4", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_cache_config.cache_config", "region", "europe-west4"),
					resource.TestCheckResourceAttr("google_vertex_ai_cache_config.cache_config", "disable_cache", "false"),
				),
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
				Config: testAccVertexAICacheConfig_region(context, "asia-northeast1", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_vertex_ai_cache_config.cache_config", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_vertex_ai_cache_config.cache_config", "region", "asia-northeast1"),
					resource.TestCheckResourceAttr("google_vertex_ai_cache_config.cache_config", "disable_cache", "true"),
				),
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

func testAccVertexAICacheConfig_region(context map[string]interface{}, region string, disableCache bool) string {
	context["region"] = region
	context["disable_cache"] = disableCache
	return acctest.Nprintf(`
resource "google_vertex_ai_cache_config" "cache_config" {
  project       = "%{project}"
  region        = "%{region}"
  disable_cache = %{disable_cache}
}
`, context)
}
