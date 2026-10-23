package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	_ "github.com/hashicorp/terraform-provider-google/google/services/apigee"
	_ "github.com/hashicorp/terraform-provider-google/google/services/compute"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
)

func TestAccApigeeSpace_basicTest(t *testing.T) {
	t.Parallel()

	context1 := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"space_id":        "test-space",
		"display_name":    "Test Space",
		"random_suffix":   acctest.RandString(t, 10),
	}

	context2 := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"space_id":        "test-space",
		"display_name":    "Updated Test Space",
		"random_suffix":   context1["random_suffix"],
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckApigeeSpaceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApigeeSpace_basicTest(context1),
			},
			{
				ResourceName:      "google_apigee_space.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApigeeSpace_basicTest(context2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_apigee_space.primary", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_apigee_space.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApigeeSpace_basicTest(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
  deletion_policy = "DELETE"
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "60s"
  depends_on = [google_project.project]
}

resource "google_project_service" "apigee" {
  project = google_project.project.project_id
  service = "apigee.googleapis.com"
  depends_on = [time_sleep.wait_60_seconds]
}

resource "google_apigee_organization" "apigee_org" {
  analytics_region    = "us-central1"
  project_id          = google_project.project.project_id
  disable_vpc_peering = true
  depends_on          = [
    google_project_service.apigee,
  ]
}

resource "google_apigee_space" "primary" {
  org_id       = google_apigee_organization.apigee_org.id
  space_id     = "%{space_id}"
  display_name = "%{display_name}"
}
`, context)
}
