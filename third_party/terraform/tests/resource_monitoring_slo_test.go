// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func setTestCheckMonitoringSloId(res string, sloId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceMonitoringSloId(res, s)
		if err != nil {
			return err
		}
		*sloId = updateId
		return nil
	}
}

func testCheckMonitoringSloIdAfterUpdate(res string, sloId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		updateId, err := getTestResourceMonitoringSloId(res, s)
		if err != nil {
			return err
		}

		if sloId == nil {
			return fmt.Errorf("unexpected error, slo ID was not set")
		}

		if *sloId != updateId {
			return fmt.Errorf("unexpected mismatch in slo ID after update, resource was recreated. Initial %q, Updated %q",
				*sloId, updateId)
		}
		return nil
	}
}

func getTestResourceMonitoringSloId(res string, s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[res]
	if !ok {
		return "", fmt.Errorf("not found: %s", res)
	}

	if rs.Primary.ID == "" {
		return "", fmt.Errorf("no ID is set for %s", res)
	}

	if v, ok := rs.Primary.Attributes["slo_id"]; ok {
		return v, nil
	}

	return "", fmt.Errorf("slo_id not set on resource %s", res)
}

func TestAccMonitoringSlo_update(t *testing.T) {
	t.Parallel()

	var generatedId string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringSloDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringSlo_basic(),
				Check:  setTestCheckMonitoringSloId("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
			{
				Config: testAccMonitoringSlo_update(),
				Check:  testCheckMonitoringSloIdAfterUpdate("google_monitoring_slo.primary", &generatedId),
			},
			{
				ResourceName:      "google_monitoring_slo.primary",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore input-only field for import
				ImportStateVerifyIgnore: []string{"service"},
			},
		},
	})
}

func testAccMonitoringSlo_basic() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.9
  rolling_period = "86400s"

  basic_sli {
    latency {
      threshold = "1s"
    }
  }
}
`
}

func testAccMonitoringSlo_update() string {
	return `
data "google_monitoring_app_engine_service" "ae" {
  module_id = "default"
}

resource "google_monitoring_slo" "primary" {
  service = data.google_monitoring_app_engine_service.ae.service_id

  goal = 0.8
  display_name = "Terraform Test updated SLO"
  calendar_period = "WEEK"

  basic_sli {
    latency {
      threshold = "2s"
    }
  }
}
`
}
