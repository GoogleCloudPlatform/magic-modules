package functions_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccProviderFunction_zone_from_id(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `
# terraform block required for provider function to be found
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

output "from_resource_id" {
  value = provider::google::zone_from_id("projects/my-project/zones/us-central1-a/disks/my-disk")
}

output "from_self_link" {
  value = provider::google::zone_from_id("https://www.googleapis.com/compute/v1/projects/my-project/zones/europe-west1-b/disks/my-disk")
}

output "from_op_resource_name" {
  value = provider::google::zone_from_id("//compute.googleapis.com/projects/my-project/zones/asia-east1-c/instances/my-instance")
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("from_resource_id", knownvalue.StringExact("us-central1-a")),
					statecheck.ExpectKnownOutputValue("from_self_link", knownvalue.StringExact("europe-west1-b")),
					statecheck.ExpectKnownOutputValue("from_op_resource_name", knownvalue.StringExact("asia-east1-c")),
				},
			},
			{
				Config: `
terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

output "invalid" {
  value = provider::google::zone_from_id("projects/my-project/regions/us-central1/subnetworks/my-subnet")
}
`,
				ExpectError: regexp.MustCompile(`doesn't\s+contain\s+the\s+expected\s+pattern\s+"zones/\{zone\}/"`),
			},
		},
	})
}
