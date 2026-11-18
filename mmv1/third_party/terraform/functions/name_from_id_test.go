package functions_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccProviderFunction_name_from_id(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `
output "from_resource_id" {
  value = provider::google::name_from_id("projects/my-project/topics/tf-test-name-id-func")
}

output "from_self_link" {
  value = provider::google::name_from_id("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a/disks/tf-test-disk-name")
}

output "from_op_resource_name" {
  value = provider::google::name_from_id("//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/my-membership")
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("from_resource_id", knownvalue.StringExact("tf-test-name-id-func")),
					statecheck.ExpectKnownOutputValue("from_self_link", knownvalue.StringExact("tf-test-disk-name")),
					statecheck.ExpectKnownOutputValue("from_op_resource_name", knownvalue.StringExact("my-membership")),
				},
			},
			{
				Config: `
output "invalid" {
  value = provider::google::name_from_id("invalid-unslashed-id")
}
`,
				ExpectError: regexp.MustCompile(`doesn't\s+contain\s+the\s+expected\s+pattern\s+"resourceType/\{name\}\$"`),
			},
		},
	})
}
