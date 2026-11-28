package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
)

func TestAccDataSourceGoogleComputeInterconnectGroups_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-interconnect-groups-ds-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInterconnectGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleComputeInterconnectGroupsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_compute_interconnect_groups.all", "id", fmt.Sprintf("projects/%s/global/interconnectGroups", envvar.GetTestProjectFromEnv())),
					// The shared test project holds other groups, so there is no count check and the list must hold this one.
					resource.TestCheckTypeSetElemNestedAttrs("data.google_compute_interconnect_groups.all", "interconnect_groups.*", map[string]string{
						"name":                         name,
						"description":                  "group for the list test",
						"intent.0.topology_capability": "NO_SLA",
					}),
					resource.TestCheckTypeSetElemAttrPair("data.google_compute_interconnect_groups.all", "interconnect_groups.*.id", "google_compute_interconnect_group.foo", "id"),
					resource.TestCheckTypeSetElemAttrPair("data.google_compute_interconnect_groups.all", "interconnect_groups.*.creation_timestamp", "google_compute_interconnect_group.foo", "creation_timestamp"),
					resource.TestCheckResourceAttrPair("data.google_compute_interconnect_groups.all", "project", "google_compute_interconnect_group.foo", "project"),
				),
			},
		},
	})
}

// testAccDataSourceGoogleComputeInterconnectGroupsConfig creates a group with no member, which needs no physical interconnect.
func testAccDataSourceGoogleComputeInterconnectGroupsConfig(name string) string {
	return fmt.Sprintf(`
resource "google_compute_interconnect_group" "foo" {
  name        = "%s"
  description = "group for the list test"
  intent {
    topology_capability = "NO_SLA"
  }
}

data "google_compute_interconnect_groups" "all" {
  project = google_compute_interconnect_group.foo.project
}
`, name)
}
