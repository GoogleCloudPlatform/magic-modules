package contactcenterinsights_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccContactCenterInsightsQaQuestion_update(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"scorecard_id":  "tf-test-qa-scorecard" + randomSuffix,
		"random_suffix": randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckContactCenterInsightsQaQuestionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContactCenterInsightsQaQuestion_update(context),
			},
			{
				ResourceName:            "google_contact_center_insights_qa_question.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "qa_scorecard", "revision"},
			},
			{
				Config: testAccContactCenterInsightsQaQuestion_update2(context),
			},
			{
				ResourceName:            "google_contact_center_insights_qa_question.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "qa_scorecard", "revision"},
			},
		},
	})
}

func testAccContactCenterInsightsQaQuestion_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_insights_qa_scorecard" "scorecard" {
  qa_scorecard_id = "%{scorecard_id}"
  location        = "us-central1"
  display_name    = "My Scorecard"
  source          = "QA_SCORECARD_SOURCE_CUSTOMER_DEFINED"
}

resource "google_contact_center_insights_qa_scorecard_revision" "rev" {
  qa_scorecard = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  location     = "us-central1"
}

resource "google_contact_center_insights_qa_question" "default" {
  qa_scorecard    = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
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
}
`, context)
}

func testAccContactCenterInsightsQaQuestion_update2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_contact_center_insights_qa_scorecard" "scorecard" {
  qa_scorecard_id = "%{scorecard_id}"
  location        = "us-central1"
  display_name    = "My Scorecard"
  source          = "QA_SCORECARD_SOURCE_CUSTOMER_DEFINED"
}

resource "google_contact_center_insights_qa_scorecard_revision" "rev" {
  qa_scorecard = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  location     = "us-central1"
}

resource "google_contact_center_insights_qa_question" "default" {
  qa_scorecard    = google_contact_center_insights_qa_scorecard.scorecard.qa_scorecard_id
  revision       = google_contact_center_insights_qa_scorecard_revision.rev.qa_scorecard_revision_id
  location       = "us-central1"

  question_body  = "Did the agent greet the customer? Updated"
  question_type  = "CUSTOMIZABLE"

  answer_choices {
    str_value = "Yes1"
    score     = 1.0
  }
  answer_choices {
    str_value = "No"
    score     = 0.25
  }
}
`, context)
}
