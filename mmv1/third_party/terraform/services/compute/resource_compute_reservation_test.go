package compute_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeReservation_resourcePolicies(t *testing.T) {
	t.Parallel()

	rand := acctest.RandString(t, 10)
	reservationName := fmt.Sprintf("tf-test-res-%s", rand)
	policyName := fmt.Sprintf("tf-test-pol-%s", rand)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeReservation_resourcePolicies(reservationName, policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"google_compute_reservation.reservation",
						"resource_policies.policy1",
						regexp.MustCompile(`resourcePolicies/`),
					),
				),
			},
			{
				ResourceName:            "google_compute_reservation.reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_policies"},
			},
		},
	})
}

func TestAccComputeReservation_update(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
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

func TestAccComputeReservation_deleteAtTime(t *testing.T) {
	acctest.SkipIfVcr(t) // timestamp
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	deleteTime := time.Now().UTC().Add(24 * time.Hour) // Set delete_at_time to 24 hours in the future
	deleteAtTimeRFC3339 := deleteTime.Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeReservation_deleteAtTime_deleteAfterDuration(reservationName, deleteAtTimeRFC3339, deleteTime.Unix()),
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				Config: testAccComputeReservation_deleteAtTime(reservationName, deleteAtTimeRFC3339),
			},
			{
				ResourceName:      "google_compute_reservation.reservation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeReservation_deleteAfterDuration(t *testing.T) {
	acctest.SkipIfVcr(t) // timestamp
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	deleteTime := time.Now().UTC().Add(24 * time.Hour) // Set delete_at_time to 24 hours in the future
	deleteAtTimeRFC3339 := deleteTime.Format(time.RFC3339)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccComputeReservation_deleteAtTime_deleteAfterDuration(reservationName, deleteAtTimeRFC3339, deleteTime.Unix()),
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
			{
				Config: testAccComputeReservation_deleteAfterDuration(reservationName, deleteTime.Unix()),
			},
			{
				ResourceName:            "google_compute_reservation.reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_after_duration"},
			},
		},
	})
}

func testAccComputeReservation_resourcePolicies(reservationName, policyName string) string {
	return fmt.Sprintf(`
resource "google_compute_resource_policy" "placement" {
  name   = "%s"
  region = "us-central1"
  // Compact policy for reservation must not set vm_count (API: incremental only).
  group_placement_policy {
    collocation = "COLLOCATED"
  }
}

resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"

  resource_policies = {
    policy1 = google_compute_resource_policy.placement.self_link
  }

  specific_reservation_required = true

  specific_reservation {
    count = 2
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }

  depends_on = [google_compute_resource_policy.placement]
}
`, policyName, reservationName)
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

func testAccComputeReservation_deleteAtTime(reservationName, time string) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"
  delete_at_time = "%s"

  specific_reservation {
    count = 2
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}
`, reservationName, time)
}

func testAccComputeReservation_deleteAfterDuration(reservationName string, duration int64) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"
  delete_after_duration {
	seconds = %d
  }

  specific_reservation {
    count = 2
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}
`, reservationName, duration)
}

func testAccComputeReservation_deleteAtTime_deleteAfterDuration(reservationName, time string, duration int64) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "reservation" {
  name = "%s"
  zone = "us-central1-a"
  delete_at_time = "%s"
  delete_after_duration {
	seconds = %d
  }

  specific_reservation {
    count = 2
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}
`, reservationName, time, duration)
}
