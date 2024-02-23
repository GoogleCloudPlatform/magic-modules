package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccProviderFunction_region_from_id(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	regionRegex := regexp.MustCompile(fmt.Sprintf("^%s$", region))

	context := map[string]interface{}{
		"function_name": "region_from_id",
		"output_name":   "region",
		"resource_name": fmt.Sprintf("tf-test-region-id-func-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Can get the region from a resource's id in one step
				// Uses google_compute_subnetwork resource's id attribute with format projects/{{project}}/regions/{{region}}/subnetworks/{{name}}
				Config: testProviderFunction_get_region_from_resource_id(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), regionRegex),
				),
			},
			{
				// Can get the region from a resource's self_link in one step
				// Uses google_compute_subnetwork resource's self_link attribute
				Config: testProviderFunction_get_region_from_resource_self_link(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), regionRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_region_from_resource_id(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
	required_providers {
		google = {
			source = "hashicorp/google"
		}
	}
}

data "google_compute_network" "default" {
  name = "default"
}

resource "google_compute_subnetwork" "default" {
  name          = "%{resource_name}"
  ip_cidr_range = "10.2.0.0/16"
  network        = data.google_compute_network.default.id
}

output "%{output_name}" {
  value = provider::google::%{function_name}(google_compute_subnetwork.default.id)
}
`, context)
}

func testProviderFunction_get_region_from_resource_self_link(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
	required_providers {
		google = {
			source = "hashicorp/google"
		}
	}
}

data "google_compute_network" "default" {
  name = "default"
}

resource "google_compute_subnetwork" "default" {
  name          = "%{resource_name}"
  ip_cidr_range = "10.2.0.0/16"
  network        = data.google_compute_network.default.id
}

output "%{output_name}" {
	value = provider::google::%{function_name}(google_compute_subnetwork.default.self_link)
}
`, context)
}
