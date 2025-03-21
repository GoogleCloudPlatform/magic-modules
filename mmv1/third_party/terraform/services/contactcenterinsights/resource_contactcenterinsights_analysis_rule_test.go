package contactcenterinsights_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestInsightsAnalysisRule(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
		"random_suffix":   acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInsightsAnalysisRule(context),
			},
			{
				ResourceName:      "google_contact_center_insights_analysis_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactCenterInsightsAnalysisRule_full(context),
			},
			{
				ResourceName:            "google_contact_center_insights_analysis_rule.basic_analysis_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccContactCenterInsightsAnalysisRule_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_contact_center_insights_analysis_rule.basic_analysis_rule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_contact_center_insights_analysis_rule.basic_analysis_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccInsightsAnalysisRule(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_project" "project" {
		name = "tf-test-insights-analysis-rule"
		project_id = "tf-test-insights-analysis-rule-%{random_suffix}"
		org_id     = "%{org_id}"
		billing_account = "%{billing_account}"
	}
	
	resource "google_contact_center_insights_analysis_rule" "default" {
	  	project = google_project.project.project_id
		name = "test-analysis-rule"
		create_time = "2024-01-01T00:00:00Z"
		update_time = "2024-01-01T00:00:00Z"
		conversation_filter = "test-filter"
		analysis_percentage = 0.5
		active = true
	}
	`, context)
}

func testAccContactCenterInsightsAnalysisRule_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_ai_insights_analysis_rule" "basic_analysis_rule" {
  name = "tf_test_basic_analysis_rule%{random_suffix}"
  display_name = "analysis-rule-display-name"
  create_time = "2025-01-01T00:00:00Z"
  update_time = "2025-01-01T00:00:00Z"
  conversation_filter = "test-filter"
  annotator_selector {
    run_interruption_annotator = true
  }
  analysis_percentage = 0.5
  active    = true
}
`, context)
}

func testAccContactCenterInsightsAnalysisRule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_ai_insights_analysis_rule" "basic_analysis_rule" {
  name = "tf_test_basic_analysis_rule%{random_suffix}"
  display_name = "analysis-rule-display-name-updated"
  create_time = "2025-01-02T00:00:00Z"
  update_time = "2025-01-02T00:00:00Z"
  conversation_filter = ""
  annotator_selector {
    run_interruption_annotator = false
  }
  analysis_percentage = 0.0
  active    = false
}
`, context)
}
