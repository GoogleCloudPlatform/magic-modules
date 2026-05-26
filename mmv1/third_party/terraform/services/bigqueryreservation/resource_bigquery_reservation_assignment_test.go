package bigqueryreservation_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/bigqueryreservation"
)

func TestAccBigqueryReservationReservationAssignment_bareNameWithoutLocation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigqueryReservationReservationAssignment_bareNameWithoutLocation(context),
				ExpectError: regexp.MustCompile("`location` is required"),
			},
		},
	})
}

func testAccBigqueryReservationReservationAssignment_bareNameWithoutLocation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "test" {
  name          = "tf-test-reservation-%{random_suffix}"
  location      = "us-central1"
  slot_capacity = 0
  edition       = "ENTERPRISE"
}

resource "google_bigquery_reservation_assignment" "test" {
  reservation = google_bigquery_reservation.test.name
  assignee    = "projects/${google_bigquery_reservation.test.project}"
  job_type    = "QUERY"
}
`, context)
}
