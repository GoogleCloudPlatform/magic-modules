package contactcenterinsights_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccContactCenterInsightsAssessmentRule_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":  acctest.RandString(t, 10),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccInsightsAssessmentRule(context),
			},
			{
				ResourceName:      "google_contact_center_insights_assessment_rule.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactCenterInsightsAssessmentRule_full(context),
			},
			{
				ResourceName:            "google_contact_center_insights_assessment_rule.basic_assessment_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"assessment_rule_id", "location"},
			},
			{
				Config: testAccContactCenterInsightsAssessmentRule_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_contact_center_insights_assessment_rule.basic_assessment_rule", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_contact_center_insights_assessment_rule.basic_assessment_rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"assessment_rule_id", "location"},
			},
		},
	})
}

func testAccInsightsAssessmentRule(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_contact_center_insights_assessment_rule" "default" {
	    display_name = "default-assessment-rule-display-name-%{random_suffix}"
		location = "us-central1"
		sample_rule {
			sample_row = 5
		}
		schedule_info {
		    schedule = "every 24 hours"
		}
		active = true
	}
	`, context)
}

func testAccContactCenterInsightsAssessmentRule_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_insights_assessment_rule" "basic_assessment_rule" {
  display_name = "assessment-rule-display-name-%{random_suffix}"
  location = "us-central1"
  conversation_filter = "agent_id = \"1\""
  sample_rule {
	sample_row = 5
  }
  schedule_info {
	schedule = "every 24 hours"
  }
  active = true
}
`, context)
}

func testAccContactCenterInsightsAssessmentRule_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_insights_assessment_rule" "basic_assessment_rule" {
  display_name = "assessment-rule-display-name-%{random_suffix}"
  location = "us-central1"
  conversation_filter = "agent_id = \"1\""
  sample_rule {
	sample_percentage = 0.5
	dimension = "quality_metadata.agent_info.agent_id"
  }
  schedule_info {
	schedule = "every 168 hours"
  }
  active = false
}
`, context)
}
