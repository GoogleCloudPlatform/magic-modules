package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeFutureReservation_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderBetaFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeFutureReservation_full(context),
			},
			{
				ResourceName:      "google_compute_future_reservation.gce_future_reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeFutureReservation_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_compute_future_reservation.gce_future_reservation", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_compute_future_reservation.gce_future_reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeFutureReservation_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_future_reservation" "gce_future_reservation" {
  provider = google-beta
  name     = "tf_test_gce_future_reservation%{random_suffix}"
  project  = "%{project}"
  zone     = "%{zone}"
  auto_created_reservations_delete_time = "2025-08-03T00:00:00Z"

  auto_created_reservations_duration {
    seconds = 86400
    nanos   = 0
  }

  auto_delete_auto_created_reservations = true

  commitment_info {
    commitment_name           = "tf-test-commitment%{random_suffix}"
    commitment_plan           = "TWELVE_MONTH"
    previous_commitment_terms = "EXTEND_TO_FIT"
  }

  deployment_type  = "MIG"
  description      = "Test future reservation description"
  name_prefix      = "tf-test-fr-px"
  planning_status  = "DRAFTING"
  reservation_name = "tf-reservation%{random_suffix}"
  scheduling_type  = "STANDARD"

  share_settings {
    share_type = "SPECIFIC_PROJECTS"
    projects   = ["%{project}"]
  }

  time_window {
    start_time = "2025-08-01T00:00:00Z"
    end_time   = "2025-08-02T00:00:00Z"
    duration {
      seconds = 86400
      nanos   = 0
    }
  }
}
`, context)
}

func testAccComputeFutureReservation_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_future_reservation" "gce_future_reservation" {
	provider = google-beta
	name     = "tf-test-fr%{random_suffix}"
	project  = "%{project}"
	zone     = "%{zone}"
  
	auto_created_reservations_delete_time = "2025-08-04T00:00:00Z"
  
	auto_created_reservations_duration {
	  seconds = 172800
	  nanos   = 0
	}
  
	auto_delete_auto_created_reservations = false
  
	commitment_info {
	  commitment_name           = "tf-test-commitment-update%{random_suffix}"
	  commitment_plan           = "ONE_YEAR"
	  previous_commitment_terms = "RETAIN"
	}
  
	deployment_type  = "MIG"
	description      = "Updated test future reservation description"
	name_prefix      = "tf-test-fr-up"
	planning_status  = "SUBMITTED"
	reservation_name = "tf-reservation-up%{random_suffix}"
	scheduling_type  = "STANDARD"
  
	share_settings {
	  share_type = "SPECIFIC_PROJECTS"
	  projects   = ["%{project}"]
	}
	time_window {
		start_time = "2025-08-01T00:00:00Z"
		end_time   = "2025-08-03T00:00:00Z"
		duration {
		  seconds = 172800
		  nanos   = 0
		}
	  }
}
`, context)
}
