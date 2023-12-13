// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeReservation(t *testing.T) {
	t.Parallel()

	reservationName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	rsName := "foobar"
	rsFullName := fmt.Sprintf("google_compute_reservation.%s", rsName)
	dsName := "my_reservation"
	dsFullName := fmt.Sprintf("data.google_compute_reservation.%s", dsName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataSourceComputeReservationDestroy(t, rsFullName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeReservationConfig(reservationName, rsName, dsName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeReservationCheck(t, dsFullName, rsFullName),
				),
			},
		},
	})
}

func testAccDataSourceComputeReservationCheck(t *testing.T, data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes

		reservation_attrs_to_test := []string{
			"name",
			"specific_reservation",
		}

		for _, attr_to_check := range reservation_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if !tpgresource.CompareSelfLinkOrResourceName("", ds_attr["self_link"], rs_attr["self_link"], nil) && ds_attr["self_link"] != rs_attr["self_link"] {
			return fmt.Errorf("self link does not match: %s vs %s", ds_attr["self_link"], rs_attr["self_link"])
		}

		if ds_attr["status"] != "READY" {
			return fmt.Errorf("status is %s; want READY", ds_attr["status"])
		}

		return nil
	}
}

func testAccCheckDataSourceComputeReservationDestroy(t *testing.T, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_reservation" {
				continue
			}

			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			_, err := config.NewComputeClient(config.UserAgent).Reservations.Get(
				config.Project, rs.Primary.Attributes["zone"], rs.Primary.Attributes["name"]).Do()
			if err == nil {
				return fmt.Errorf("Reservation still exists")
			}
		}

		return nil
	}
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
