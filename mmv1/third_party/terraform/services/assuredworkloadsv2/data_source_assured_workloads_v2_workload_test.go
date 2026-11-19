package assuredworkloadsv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/assuredworkloadsv2"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

func TestAccDataSourceAssuredWorkloadsV2Workload_basic(t *testing.T) {
	t.Parallel()

	context := testAccAssuredWorkloadsV2Context(t)

	ignoredFields := []string{
		"resource_config.#",
		"resource_config.0.%",
		"resource_config.0.existing_project",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAssuredWorkloadsV2WorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssuredWorkloadsV2Workload_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_assured_workloads_v2_workload.by_ids",
						"google_assured_workloads_v2_workload.workload",
						ignoredFields,
					),
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_assured_workloads_v2_workload.by_name",
						"google_assured_workloads_v2_workload.workload",
						ignoredFields,
					),
				),
			},
		},
	})
}

func testAccDataSourceAssuredWorkloadsV2Workload_basic(context map[string]interface{}) string {
	return testAccAssuredWorkloadsV2Workload_existingProjectMonitorRev62(context) + `
data "google_assured_workloads_v2_workload" "by_ids" {
  organization = google_assured_workloads_v2_workload.workload.organization
  location     = google_assured_workloads_v2_workload.workload.location
  workload_id  = google_assured_workloads_v2_workload.workload.workload_id
}

data "google_assured_workloads_v2_workload" "by_name" {
  name = google_assured_workloads_v2_workload.workload.name
}
`
}
