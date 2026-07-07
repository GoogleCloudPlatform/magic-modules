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
