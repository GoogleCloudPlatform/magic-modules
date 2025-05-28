package compute_test

import (
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeSnapshotSettings_snapshotSettings_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSnapshotSettings_snapshotSettings_basic(context),
			},
			{
				ResourceName:      "google_compute_snapshot_settings.tf_test_snapshot_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeSnapshotSettings_snapshotSettings_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_snapshot_settings.tf_test_snapshot_settings", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_snapshot_settings.tf_test_snapshot_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeSnapshotSettings_snapshotSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_snapshot_settings" "tf_test_snapshot_settings" {
    project   = "%{project}"
    storage_location {
        policy    = "SPECIFIC_LOCATIONS"
        locations {
            name     = "us-central1"
            location = "us-central1"
        }
    }
}
`, context)
}

func testAccComputeSnapshotSettings_snapshotSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_snapshot_settings" "tf_test_snapshot_settings" {
    project   = "%{project}"
    storage_location {
        policy    = "NEAREST_MULTI_REGION"
    }
}
`, context)
}
