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

func TestAccProviderFunction_region_from_zone(t *testing.T) {
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

output "standard_zone" {
  value = provider::google::region_from_zone("us-central1-a")
}

output "multi_hyphen_zone" {
  value = provider::google::region_from_zone("northamerica-northeast1-b")
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("standard_zone", knownvalue.StringExact("us-central1")),
					statecheck.ExpectKnownOutputValue("multi_hyphen_zone", knownvalue.StringExact("northamerica-northeast1")),
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

output "invalid_zone" {
  value = provider::google::region_from_zone("invalidzone")
}
`,
				ExpectError: regexp.MustCompile(`is\s+not\s+a\s+valid\s+zone\s+name`),
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

output "empty_zone" {
  value = provider::google::region_from_zone("")
}
`,
				ExpectError: regexp.MustCompile(`cannot\s+be\s+empty`),
			},
		},
	})
}
