// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccDataSourceGoogleServiceAttachment(t *testing.T) {
	t.Parallel()
	
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleServiceAttachmentConfig(context),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleServiceAttachmentCheck("data.google_compute_service_attachment.my_attachment", "google_compute_service_attachment.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleServiceAttachmentCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
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
		attachment_attrs_to_test := []string{
			"id",
			"name",
			"region",
		}

		for _, attr_to_check := range attachment_attrs_to_test {
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

		return nil
	}
}

func testAccDataSourceGoogleServiceAttachmentConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_service_attachment" "foobar" {
  name                     = "tf-test%{random_suffix}"
  description              = "my-description"
  region                   = "us-west1"
  enable_proxy_protocol    = true
  connection_preference    = "ACCEPT_AUTOMATIC"
  nat_subnets              = [google_compute_subnetwork.psc_ilb_nat.id]
  target_service           = google_compute_forwarding_rule.psc_ilb_target_service.id
}

resource "google_compute_forwarding_rule" "psc_ilb_target_service" {
	name   = "producer-forwarding-rule%{random_suffix}"
	region = "us-west1"
  
	load_balancing_scheme = "INTERNAL"
	backend_service       = google_compute_region_backend_service.producer_service_backend.id
	all_ports             = true
	network               = google_compute_network.psc_ilb_network.name
	subnetwork            = google_compute_subnetwork.psc_ilb_producer_subnetwork.name
  }

  resource "google_compute_region_backend_service" "producer_service_backend" {
	name   = "producer-service"
	region = "us-west1"
  
	health_checks = [google_compute_health_check.producer_service_health_check.id]
  }
  
  resource "google_compute_health_check" "producer_service_health_check" {
	name = "producer-service-health-check%{random_suffix}"
  
	check_interval_sec = 1
	timeout_sec        = 1
	tcp_health_check {
	  port = "80"
	}
  }
  resource "google_compute_network" "psc_ilb_network" {
	name = "psc-ilb-network%{random_suffix}"
	auto_create_subnetworks = false
  }
  
  resource "google_compute_subnetwork" "psc_ilb_producer_subnetwork" {
	name   = "psc-ilb-producer-subnetwork%{random_suffix}"
	region = "us-west1"
  
	network       = google_compute_network.psc_ilb_network.id
	ip_cidr_range = "10.0.0.0/16"
  }
  
  resource "google_compute_subnetwork" "psc_ilb_nat" {
	name   = "psc-ilb-nat%{random_suffix}"
	region = "us-west1"
  
	network       = google_compute_network.psc_ilb_network.id
	purpose       =  "PRIVATE_SERVICE_CONNECT"
	ip_cidr_range = "10.1.0.0/16"
  }
data "google_compute_service_attachment" "my_attachment" {
  name = google_compute_service_attachment.foobar.name
  region = google_compute_service_attachment.foobar.region
}
`, context)
}