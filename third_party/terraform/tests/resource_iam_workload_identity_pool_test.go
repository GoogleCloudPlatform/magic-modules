package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIAMWorkloadIdentityPool_example(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"org_domain":    getTestOrgDomainFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMWorkloadIdentityPool_example(context),
			},
			{
				ResourceName:            "google_iam_workload_identity_pool.my_pool",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func testAccIAMWorkloadIdentityPool_example(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "my_project" {
  project_id = "tf-test%{random_suffix}"
  name       = "tf-test%{random_suffix}"
  org_id     = "%{org_id}"
}

resource "google_project_service" "my_service" {
  project = google_project.my_project.project_id
  service = "iam.googleapis.com"
}

resource "google_iam_workload_identity_pool" "my_pool" {
  project      = google_project_service.my_service.project
  display_name = "Name of pool"
  description  = "Identity pool for automated test"
  disabled     = true
}
`, context)
}
