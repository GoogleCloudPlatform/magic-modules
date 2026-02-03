package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccComputeRolloutPlan_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	resourceName := "google_compute_rollout_plan.acceptance"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeRolloutPlan_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-rollout-plan-%s", context["random_suffix"])),
					resource.TestCheckResourceAttr(resourceName, "description", "A test rollout plan"),
					resource.TestCheckResourceAttr(resourceName, "location_scope", "ZONAL"),
					resource.TestCheckResourceAttr(resourceName, "waves.0.display_name", "wave-1"),
					resource.TestCheckResourceAttr(resourceName, "waves.0.selectors.0.location_selector.0.included_locations.0", "us-central1-a"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeRolloutPlan_basic(context map[string]interface{}) string {
	return fmt.Sprintf(`
resource "google_compute_rollout_plan" "acceptance" {
  name           = "tf-test-rollout-plan-%s"
  description    = "A test rollout plan"
  location_scope = "ZONAL"

  waves {
    display_name = "wave-1"
    selectors {
      location_selector {
        included_locations = ["us-central1-a"]
      }
    }
    validation {
      type = "manual"
    }
  }
}
`, context["random_suffix"])
}
