// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccComputeInstanceListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	zone := envvar.GetTestZoneFromEnv()
	name := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceListBasic(zone, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_compute_instance.test", "zone", zone),
					resource.TestCheckResourceAttr("google_compute_instance.test", "project", project),
					resource.TestCheckResourceAttr("google_compute_instance.test", "name", name),
				),
			},
			{
				Query:  true,
				Config: testAccComputeInstanceListQuery(project, zone),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_compute_instance.all_in_zone", map[string]knownvalue.Check{
						"zone":    knownvalue.StringExact(zone),
						"project": knownvalue.StringExact(project),
						"name":    knownvalue.StringExact(name),
					}),
					querycheck.ExpectLengthAtLeast("google_compute_instance.all_in_zone", 1),
				},
			},
		},
	})
}

func testAccComputeInstanceListBasic(zone, name string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test" {
  name         = %q
  zone         = %q
  machine_type = "e2-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"

    access_config {
    }
  }
}
`, name, zone)
}

func testAccComputeInstanceListQuery(project, zone string) string {
	return fmt.Sprintf(`
list "google_compute_instance" "all_in_zone" {
  provider = google

  config {
    project = %q
    zone    = %q
  }
}
`, project, zone)
}
