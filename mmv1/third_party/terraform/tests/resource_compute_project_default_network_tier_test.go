package google_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeProjectDefaultNetworkTier_basic(t *testing.T) {
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billingId := acctest.GetTestBillingAccountFromEnv(t)
	projectID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:  func() { acctest.TestAccPreCheck(t) },
		Providers: acctest.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_defaultNetworkTier_premium(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_default_network_tier.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeProjectDefaultNetworkTier_modify(t *testing.T) {
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billingId := acctest.GetTestBillingAccountFromEnv(t)
	projectID := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:  func() { acctest.TestAccPreCheck(t) },
		Providers: acctest.TestAccProvidersrovidersroviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeProject_defaultNetworkTier_premium(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_default_network_tier.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccComputeProject_defaultNetworkTier_standard(projectID, pname, org, billingId),
			},
			{
				ResourceName:      "google_compute_project_default_network_tier.fizzbuzz",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeProject_defaultNetworkTier_premium(projectID, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_project_default_network_tier" "fizzbuzz" {
  project      = google_project.project.project_id
  network_tier = "PREMIUM"
  depends_on   = [google_project_service.compute]
}
`, projectID, name, org, billing)
}

func testAccComputeProject_defaultNetworkTier_standard(projectID, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_compute_project_default_network_tier" "fizzbuzz" {
  project      = google_project.project.project_id
  network_tier = "STANDARD"
  depends_on   = [google_project_service.compute]
}
`, projectID, name, org, billing)
}
