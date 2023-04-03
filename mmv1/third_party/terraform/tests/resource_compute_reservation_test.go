package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeReservation_update(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeReservation_basic(reservationName, "2"),
			},
			{
				ResourceName:      "google_compute_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccComputeReservation_basic(reservationName, "1"),
			},
			{
				ResourceName:      "google_compute_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeReservation_basic(reservationName, count string) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"

  specific_reservation {
    count = %s
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}
`, reservationName, count)
}

func TestAccComputeReservation_projectMapUpdate(t *testing.T) {
	SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"project":         GetTestProjectFromEnv(),
		"org_id":          GetTestOrgFromEnv(t),
		"billing_account": GetTestBillingAccountFromEnv(t),
		"random_suffix":   RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeReservation_projectMap(context),
			},
			{
				ResourceName:            "google_compute_reservation.gce_reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccComputeReservation_multipleProjectMap(context),
			},
			{
				ResourceName:            "google_compute_reservation.gce_reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
			{
				Config: testAccComputeReservation_projectMap(context),
			},
			{
				ResourceName:            "google_compute_reservation.gce_reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone"},
			},
		},
	})
}

func testAccComputeReservation_projectMap(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "owner_project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}


resource "google_project_service" "compute" {
  project = google_project.owner_project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project" "guest_project" {
  project_id      = "tf-test-2%{random_suffix}"
  name            = "tf-test-2%{random_suffix}"
  org_id          = "%{org_id}"
}

resource "google_organization_policy" "shared_reservation_org_policy" {
  org_id     = "%{org_id}"
  constraint = "constraints/compute.sharedReservationsOwnerProjects"
  list_policy {
    allow {
      values = ["projects/${google_project.owner_project.number}"]
    }
  }
}

resource "google_compute_reservation" "gce_reservation" {
  project = google_project.owner_project.project_id
  name = "tf-test-gce-shared-reservation%{random_suffix}"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
  share_settings {
    share_type = "SPECIFIC_PROJECTS"
    project_map {
      id = google_project.guest_project.number
      project_id = google_project.guest_project.number
    }
  }
  depends_on = [google_organization_policy.shared_reservation_org_policy,google_project_service.compute]
}
`, context)
}

func testAccComputeReservation_multipleProjectMap(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "owner_project" {
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "compute" {
  project = google_project.owner_project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project" "guest_project" {
  project_id      = "tf-test-2%{random_suffix}"
  name            = "tf-test-2%{random_suffix}"
  org_id          = "%{org_id}"
}

resource "google_project" "guest_project_1" {
  project_id      = "tf-test-3%{random_suffix}"
  name            = "tf-test-3%{random_suffix}"
  org_id          = "%{org_id}"
}

resource "google_organization_policy" "shared_reservation_org_policy" {
  org_id     = "%{org_id}"
  constraint = "constraints/compute.sharedReservationsOwnerProjects"
  list_policy {
    allow {
      values = ["projects/${google_project.owner_project.number}"]
    }
  }
}

resource "google_compute_reservation" "gce_reservation" {
  project = google_project.owner_project.project_id
  name = "tf-test-gce-shared-reservation%{random_suffix}"
  zone = "us-central1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
  share_settings {
    share_type = "SPECIFIC_PROJECTS"
    project_map {
      id = google_project.guest_project.number
      project_id = google_project.guest_project.number
    }

    project_map {
      id = google_project.guest_project_1.number
      project_id = google_project.guest_project_1.number
    }
  }
  depends_on = [google_organization_policy.shared_reservation_org_policy,google_project_service.compute]
}
`, context)
}
