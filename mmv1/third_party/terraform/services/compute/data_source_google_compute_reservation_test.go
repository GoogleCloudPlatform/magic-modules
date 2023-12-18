// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestComputeReservationIdParsing(t *testing.T) {
	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedCanonicalId string
		Config              *transport_tpg.Config
	}{
		"id is a full self link": {
			ImportId:            "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/reservations/test-reservation",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/zones/us-central1-a/reservations/test-reservation",
		},
		"id is a partial self link": {
			ImportId:            "projects/test-project/zones/us-central1-a/reservations/test-reservation",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/zones/us-central1-a/reservations/test-reservation",
		},
		"id is project/region/address": {
			ImportId:            "test-project/us-central1-a/test-reservation",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/test-project/zones/us-central1-a/reservations/test-reservation",
		},
		"id is region/address": {
			ImportId:            "us-central1-a/test-reservation",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/zones/us-central1-a/reservations/test-reservation",
			Config:              &transport_tpg.Config{Project: "default-project"},
		},
		"id is address": {
			ImportId:            "test-reservation",
			ExpectedError:       false,
			ExpectedCanonicalId: "projects/default-project/zones/us-east1-a/reservations/test-reservation",
			Config:              &transport_tpg.Config{Project: "default-project", Zone: "us-east1-a"},
		},
		"id has invalid format": {
			ImportId:      "i/n/v/a/l/i/d",
			ExpectedError: true,
		},
	}

	for tn, tc := range cases {
		addressId, err := compute.ParseComputeReservationId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if addressId.CanonicalId() != tc.ExpectedCanonicalId {
			t.Fatalf("bad: %s, expected canonical id to be `%s` but is `%s`", tn, tc.ExpectedCanonicalId, addressId.CanonicalId())
		}
	}
}

func TestAccDataSourceComputeReservation(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	rsName := "foobar"
	dsName := "my_reservation"
	rsFullName := fmt.Sprintf("google_compute_reservation.%s", rsName)
	dsFullName := fmt.Sprintf("data.google_compute_reservation.%s", dsName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeReservationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeReservationConfig(reservationName, rsName, dsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsFullName, "status", "READY"),
					acctest.CheckDataSourceStateMatchesResourceState(dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeReservationConfig(reservationName, rsName, dsName string) string {
	return fmt.Sprintf(`
resource "google_compute_reservation" "%s" {
  name = "%s"
  zone = "us-west1-a"

  specific_reservation {
    count = 1
    instance_properties {
      min_cpu_platform = "Intel Cascade Lake"
      machine_type     = "n2-standard-2"
    }
  }
}

data "google_compute_reservation" "%s" {
  name = google_compute_reservation.%s.name
  zone = "us-west1-a"
}
`, rsName, reservationName, dsName, rsName)
}
