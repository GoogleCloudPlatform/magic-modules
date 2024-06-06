package bigqueryreservation_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryReservationAssignment_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservationAssignment_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_bigquery_reservation_assignment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reservation"},
			},
		},
	})
}

func testAccBigqueryReservationAssignment_BasicHandWritten(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_reservation" "basic" {
  name  = "tf-test-my-reservation%{random_suffix}"
  project = "%{project_name}"
  location = "us-central1"
  slot_capacity = 0
  ignore_idle_slots = false
}

resource "google_bigquery_reservation_assignment" "primary" {
  assignee  = "projects/%{project_name}"
  job_type = "PIPELINE"
  reservation = google_bigquery_reservation.basic.id
}
`, context)
}
