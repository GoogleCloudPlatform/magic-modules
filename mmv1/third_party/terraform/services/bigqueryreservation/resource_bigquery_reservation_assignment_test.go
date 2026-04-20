package bigqueryreservation_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigQueryReservationAssignment_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryReservationAssignment_basic(context),
			},
			{
				ResourceName:      "google_bigquery_reservation_assignment.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigQueryReservationAssignment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "basic" {
  name            = "tf-test-res-%{random_suffix}"
  location        = "us-central1"
  slot_capacity   = 100
  ignore_idle_slots = false
}

resource "google_bigquery_reservation_assignment" "primary" {
  assignee    = "projects/%{project_id}"
  job_type    = "QUERY"
  principal   = "principal://iam.googleapis.com/projects/-/serviceAccounts/%{service_account}"
  location    = "us-central1"
  reservation = google_bigquery_reservation.basic.id
}
`, map[string]interface{}{
		"random_suffix":  context["random_suffix"],
		"project_id":     envvar.GetTestProjectFromEnv(),
		"service_account": envvar.GetTestServiceAccountFromEnv(),
	})
}