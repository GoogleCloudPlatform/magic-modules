package contactcenterinsights_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContactCenterInsightsQaQuestionTag_full(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"qa_question_tag_id": "tf-test-tag-" + randomSuffix,
		"random_suffix":      randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContactCenterInsightsQaQuestionTag_step1(context),
			},
			{
				Config: testAccContactCenterInsightsQaQuestionTag_step2(context),
			},
		},
	})
}

func testAccContactCenterInsightsQaQuestionTag_step1(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_contact_center_insights_qa_scorecard" "scorecard" {
  qa_scorecard_id = "tf-test-%{random_suffix}"
  location        = "us-central1"
  display_name    = "My Scorecard"
  source          = "QA_SCORECARD_SOURCE_CUSTOMER_DEFINED"
}

resource "google_contact_center_insights_qa_scorecard_revision" "rev" {
  qa_scorecard = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  location     = "us-central1"
}

resource "google_contact_center_insights_qa_question" "question" {
  qa_scorecard   = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  revision       = google_contact_center_insights_qa_scorecard_revision.rev.qa_scorecard_revision_id
  location       = "us-central1"

  question_body  = "Did the agent greet the customer?"
  question_type  = "CUSTOMIZABLE"

  answer_choices {
    str_value = "Yes"
    score     = 1.0
  }
  answer_choices {
    str_value = "No"
    score     = 0.5
  }

  tags = ["projects/${data.google_project.project.number}/locations/us-central1/qaQuestionTags/%{qa_question_tag_id}"]
}

resource "google_contact_center_insights_qa_question_tag" "default" {
  qa_question_tag_id = "%{qa_question_tag_id}"
  location           = "us-central1"
  display_name       = "My Question Tag %{qa_question_tag_id}"
  qa_question_ids    = ["projects/${data.google_project.project.number}/locations/us-central1/qaScorecards/${google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id}/revisions/${google_contact_center_insights_qa_scorecard_revision.rev.qa_scorecard_revision_id}/qaQuestions/${google_contact_center_insights_qa_question.question.name}"]
}
`, context)
}

func testAccContactCenterInsightsQaQuestionTag_step2(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_contact_center_insights_qa_scorecard" "scorecard" {
  qa_scorecard_id = "tf-test-%{random_suffix}"
  location        = "us-central1"
  display_name    = "My Scorecard"
  source          = "QA_SCORECARD_SOURCE_CUSTOMER_DEFINED"
}

resource "google_contact_center_insights_qa_scorecard_revision" "rev" {
  qa_scorecard = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  location     = "us-central1"
}

resource "google_contact_center_insights_qa_question" "question" {
  qa_scorecard   = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  revision       = google_contact_center_insights_qa_scorecard_revision.rev.qa_scorecard_revision_id
  location       = "us-central1"

  question_body  = "Did the agent greet the customer?"
  question_type  = "CUSTOMIZABLE"

  answer_choices {
    str_value = "Yes"
    score     = 1.0
  }
  answer_choices {
    str_value = "No"
    score     = 0.5
  }
  
  tags = []
}

resource "google_contact_center_insights_qa_question_tag" "default" {
  qa_question_tag_id = "%{qa_question_tag_id}"
  location           = "us-central1"
  display_name       = "My Question Tag %{qa_question_tag_id}"
  qa_question_ids    = []
}
`, context)
}
