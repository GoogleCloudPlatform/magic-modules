package bigqueryreservation_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
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

func TestAccBigqueryReservationReservationAssignment_customSymmetricReservationDiffSuppress(t *testing.T) {
	t.Parallel()

	randomSuffix := acctest.RandString(t, 10)

	context := map[string]interface{}{
		"project":          envvar.GetTestProjectFromEnv(),
		"reservation_name": "tf-test-example-res-" + randomSuffix,
		"random_suffix":    randomSuffix,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigqueryReservationReservationAssignmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservationReservationAssignment_customSymmetricReservationDiffSuppress(context),
			},
			{
				ResourceName:            "google_bigquery_reservation_assignment.assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"}, // Verify reservation is NOT ignored and passes verification!
			},
			{
				ResourceName:       "google_bigquery_reservation_assignment.assignment",
				RefreshState:       true,
				ExpectNonEmptyPlan: false, // Verify plan is clean post-import!
				ImportStateKind:    resource.ImportBlockWithResourceIdentity,
			},
		},
	})
}

func testAccBigqueryReservationReservationAssignment_customSymmetricReservationDiffSuppress(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "basic" {
  name  = "%{reservation_name}"
  project = "%{project}"
  location = "us-central1"
  slot_capacity = 0
  ignore_idle_slots = false
}

resource "google_bigquery_reservation_assignment" "assignment" {
  assignee  = "projects/%{project}"
  job_type = "PIPELINE"
  reservation = google_bigquery_reservation.basic.id
  location = "us-central1"
}
`, context)
}
