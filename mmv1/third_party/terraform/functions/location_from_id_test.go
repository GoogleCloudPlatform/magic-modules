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

func TestAccProviderFunction_location_from_id(t *testing.T) {
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

output "from_project_location_id" {
  value = provider::google::location_from_id("projects/my-project/locations/us-central1/services/my-service")
}

output "from_cloud_run_v1_id" {
  value = provider::google::location_from_id("locations/us-east1/namespaces/my-project/services/my-service")
}

output "from_self_link" {
  value = provider::google::location_from_id("https://run.googleapis.com/v2/projects/my-project/locations/europe-west1/services/my-service")
}

output "from_op_resource_name" {
  value = provider::google::location_from_id("//gkehub.googleapis.com/projects/my-project/locations/asia-east1/memberships/my-membership")
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("from_project_location_id", knownvalue.StringExact("us-central1")),
					statecheck.ExpectKnownOutputValue("from_cloud_run_v1_id", knownvalue.StringExact("us-east1")),
					statecheck.ExpectKnownOutputValue("from_self_link", knownvalue.StringExact("europe-west1")),
					statecheck.ExpectKnownOutputValue("from_op_resource_name", knownvalue.StringExact("asia-east1")),
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
  value = provider::google::location_from_id("projects/my-project/zones/us-central1-a/instances/my-instance")
}
`,
				ExpectError: regexp.MustCompile(`doesn't\s+contain\s+the\s+expected\s+pattern\s+"locations/\{location\}/"`),
			},
		},
	})
}
