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

	context := map[string]interface{}{
		"function_name": "region_from_zone",
		"output_name":   "zone",
		"resource_name": fmt.Sprintf("tf-test-region-from-zone-func-%s", acctest.RandString(t, 10)),
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testProviderFunction_get_region_from_zone(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), "us-central1"),
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

resource "google_filestore_instance" "instance" {
	name = "%{resource_name}"
	location = "us-central1-b"
	tier     = "BASIC_HDD"
  
	file_shares {
	  capacity_gb = 1024
	  name        = "share1"
	}
  
	networks {
	  network = "default"
	  modes   = ["MODE_IPV4"]
	}
  }

output "%{output_name}" {
	value = provider::google::%{function_name}(google_filestore_instance.default.location)
}
`, context)
}
