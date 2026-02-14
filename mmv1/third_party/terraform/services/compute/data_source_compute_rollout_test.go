package compute_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"testing"
)

func TestAccDataSourceComputeRollouts_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"pid": envvar.GetTestProjectFromEnv(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeRollouts_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_rollouts.all", "project", context["pid"].(string)),
					resource.TestCheckResourceAttrSet("data.google_compute_rollouts.all", "id"),
				),
			},
		},
	})
}

func testAccDataSourceComputeRollouts_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_rollouts" "all" {
  project = "%{pid}"
}
`, context)
}
