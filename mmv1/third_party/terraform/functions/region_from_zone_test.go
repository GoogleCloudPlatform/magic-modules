// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccProviderFunction_region_from_zone(t *testing.T) {
	t.Parallel()

	projectZone := envvar.GetTestZoneFromEnv()
	projectZoneRegex := regexp.MustCompile(fmt.Sprintf("^%s$", projectZone[:len(projectZone)-2]))

	context := map[string]interface{}{
		"function_name":     "region_from_zone",
		"output_name":       "zone",
		"resource_name":     fmt.Sprintf("tf-test-region-from-zone-func-%s", acctest.RandString(t, 10)),
		"resource_location": projectZone,
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testProviderFunction_get_region_from_zone(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), projectZoneRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_region_from_zone(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
	required_providers {
		google = {
			source = "hashicorp/google"
		}
	}
}

resource "google_cloud_run_service" "default" {
	name     = "%{resource_name}"
	location = "%{resource_location}"
  
	template {
	  spec {
		containers {
		  image = "us-docker.pkg.dev/cloudrun/container/hello"
		}
	  }
	}
  
	traffic {
	  percent         = 100
	  latest_revision = true
	}
  }

output "%{output_name}" {
	value = provider::google::%{function_name}(google_cloud_run_service.default.location)
}
`, context)
}
