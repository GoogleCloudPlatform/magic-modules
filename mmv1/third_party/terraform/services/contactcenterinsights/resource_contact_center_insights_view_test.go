package contactcenterinsights_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestInsightsView(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":    envvar.GetTestProjectFromEnv(),
		"region":          "us-central1",
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInsightsView(context),
			},
			{
				ResourceName:      "google_contact_center_insights_view.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactCenterInsightsView_full(context),
			},
			{
				ResourceName:            "google_contact_center_insights_view.full_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContactCenterInsightsView_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_contact_center_insights_view.full_view", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_contact_center_insights_view.full_view",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccInsightsView(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "project" {
		name = "tf-test-insights-view"
		project_id = "tf-test-insights-view-%{random_suffix}"
		org_id     = "%{org_id}"
		billing_account = "%{billing_account}"
	}
	
	resource "google_contact_center_insights_view" "default" {
	  	project = google_project.project.project_id
		location = "%{region}"
		display_name = "test-view"
		create_time = "2024-01-01T00:00:00Z"
		update_time = "2024-01-01T00:00:00Z"
		value    = "medium=\"PHONE_CALL\""
	}
	`, context)
}

func testAccContactCenterInsightsView_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_insights_view" "full_view" {
  project = "%{project_name}"
  name = "tf-insights-view-{%random_suffix}"
  location = "%{region}"
  display_name = "view-display-name-%{random_suffix}"
  create_time = "2025-01-01T00:00:00Z"
  update_time = "2025-01-01T00:00:00Z"
  value    = "medium=\"PHONE_CALL\""
}
`, context)
}

func testAccContactCenterInsightsView_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_insights_view" "full_view" {
  project = "%{project_name}"
  name = "tf-insights-view-{%random_suffix}"
  location = "%{region}"
  display_name = "view-display-name-%{random_suffix}-updated"
  create_time = "2025-01-02T00:00:00Z"
  update_time = "2025-01-02T00:00:00Z"
  value    = "medium=\"CHAT\""
}
`, context)
}
