package functions_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccProviderFunction_project_from_id(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `
output "from_resource_id" {
  value = provider::google::project_from_id("projects/my-project/topics/my-topic")
}

output "from_self_link" {
  value = provider::google::project_from_id("https://bigquery.googleapis.com/bigquery/v2/projects/my-bq-project/datasets/my_dataset")
}

output "from_op_resource_name" {
  value = provider::google::project_from_id("//gkehub.googleapis.com/projects/my-op-project/locations/us-central1/memberships/my-membership")
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("from_resource_id", knownvalue.StringExact("my-project")),
					statecheck.ExpectKnownOutputValue("from_self_link", knownvalue.StringExact("my-bq-project")),
					statecheck.ExpectKnownOutputValue("from_op_resource_name", knownvalue.StringExact("my-op-project")),
				},
			},
			{
				Config: `
output "invalid" {
  value = provider::google::project_from_id("organizations/123456789/locations/global")
}
`,
				ExpectError: regexp.MustCompile(`doesn't\s+contain\s+the\s+expected\s+pattern\s+"projects/\{project\}/"`),
			},
		},
	})
}
