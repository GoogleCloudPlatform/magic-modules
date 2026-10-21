package workbench_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/workbench"
)

func TestAccDataSourceWorkbenchInstance_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"location":      envvar.GetTestZoneFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckWorkbenchInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkbenchInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					// Ignore client-side virtual field desired_state and input-only fields (instance_owners, gce_setup.0.boot_disk.0.disk_type, gce_setup.0.data_disks.0.disk_type, gce_setup.0.vm_image) not returned by the API.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_workbench_instance.default",
						"google_workbench_instance.instance",
						[]string{
							"desired_state",
							"instance_owners",
							"gce_setup.0.boot_disk.0.disk_type",
							"gce_setup.0.data_disks.0.disk_type",
							"gce_setup.0.vm_image",
						},
					),
					// Ignore client-side virtual field desired_state and input-only fields (instance_owners, gce_setup.0.boot_disk.0.disk_type, gce_setup.0.data_disks.0.disk_type, gce_setup.0.vm_image) not returned by the API.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_workbench_instance.default_no_optional",
						"google_workbench_instance.instance",
						[]string{
							"desired_state",
							"instance_owners",
							"gce_setup.0.boot_disk.0.disk_type",
							"gce_setup.0.data_disks.0.disk_type",
							"gce_setup.0.vm_image",
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceWorkbenchInstance_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_workbench_instance" "instance" {
  name     = "tf-test-workbench-instance-%{random_suffix}"
  location = "%{location}"

  gce_setup {
    machine_type = "n4-standard-2"
    boot_disk {
      disk_type = "HYPERDISK_BALANCED"
    }
  }

  labels = {
    k = "val"
  }
}

data "google_workbench_instance" "default" {
  name     = google_workbench_instance.instance.name
  location = google_workbench_instance.instance.location
  project  = google_workbench_instance.instance.project
}

data "google_workbench_instance" "default_no_optional" {
  name = google_workbench_instance.instance.name

  depends_on = [
    google_workbench_instance.instance,
  ]
}
`, context)
}
