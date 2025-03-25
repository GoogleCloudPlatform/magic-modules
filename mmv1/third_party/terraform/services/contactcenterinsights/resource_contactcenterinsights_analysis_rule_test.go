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
		location = "us-central1"
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
  location = "us-central1"
  create_time = "2025-01-01T00:00:00Z"
  update_time = "2025-01-01T00:00:00Z"
  conversation_filter = "test-filter"
  annotator_selector {
    run_interruption_annotator = true
	issue_models    = "some_issue_model_id"
    phrase_matchers = "some_phrase_matcher_id"
    qa_config {
      scorecard_list {
        qa_scorecard_revisions = "some_scorecard_revision_id"
      }
    }
    run_entity_annotator         = true
    run_intent_annotator         = true
    run_issue_model_annotator    = true
    run_phrase_matcher_annotator = true
    run_qa_annotator             = true
    run_sentiment_annotator      = true
    run_silence_annotator        = true
    run_summarization_annotator  = true
    summarization_config {
      conversation_profile = "some_conversation_profile"
      summarization_model  = BASELINE_MODEL
    }
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
  location = "us-central1"
  create_time = "2025-01-02T00:00:00Z"
  update_time = "2025-01-02T00:00:00Z"
  conversation_filter = ""
  annotator_selector {
    run_interruption_annotator = false
	issue_models    = "alt_issue_model_id"
    phrase_matchers = "alt_phrase_matcher_id"
    qa_config {
      scorecard_list {
        qa_scorecard_revisions = "alt_scorecard_revision_id"
      }
    }
    run_entity_annotator         = false
    run_intent_annotator         = false
    run_issue_model_annotator    = false
    run_phrase_matcher_annotator = false
    run_qa_annotator             = false
    run_sentiment_annotator      = false
    run_silence_annotator        = false
    run_summarization_annotator  = false
    summarization_config {
      conversation_profile = "alt_conversation_profile"
      summarization_model  = BASELINE_MODEL_V2_0
    }
  }
  analysis_percentage = 0.0
  active    = false
}
`, context)
}
