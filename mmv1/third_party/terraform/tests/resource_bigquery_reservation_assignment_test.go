package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigqueryReservationReservationAssignment_bigqueryReservationAssignmentBasic(t *testing.T) {
	t.Parallel()

	location := "asia-northeast1"
	context := map[string]interface{}{
		"random_suffix":   randString(t, 10),
		"location":        location,
		"project_id":      fmt.Sprintf("tf-test-%d", randInt(t)),
		"project_name":    pname,
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigqueryReservationReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservationReservationAssignment_bigqueryReservationAssignmentBasic(context),
			},
			{
				ResourceName:      "google_bigquery_reservation_assignment.reservation_assignment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigqueryReservationReservationAssignment_bigqueryReservationAssignmentBasic(context map[string]interface{}) string {
	return Nprintf(`
// Create a separate project because assignment is unique per assignee and job_type
// and there are only three valid job types.
resource "google_project" "project" {
  project_id      = "%{project_id}"
  name            = "%{project_name}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "bigqueryreservation" {
  project = google_project.project.project_id
  service = "bigqueryreservation.googleapis.com"
}

resource "google_bigquery_reservation" "reservation" {
	name           = "tf-test-reservation%{random_suffix}"
	project        = google_project_service.bigqueryreservation.project
	location       = "%{location}"
	// Set to 0 for testing purposes
	// In reality this would be larger than zero
	slot_capacity  = 0
	ignore_idle_slots = false
}

resource "google_bigquery_reservation_assignment" "reservation_assignment" {
	reservation = google_bigquery_reservation.reservation.id
    assignee    = "projects/%{project_id}"
    job_type    = "QUERY"
}
`, context)
}
