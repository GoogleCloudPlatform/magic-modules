package apigee_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccApigeeSecurityActionListQuery exercises the list resource for
// google_apigee_security_action. It does not create any resources: it
// issues a pure list query against an environment that already exists in the
// test org, which is sufficient to validate that the list endpoint is reachable
// and the provider plumbing is correct.
func TestAccApigeeSecurityActionListQuery(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	orgId := envvar.GetTestOrgFromEnv(t)
	// env_id is just the environment name, not the full path
	envId := "my-test-environment"

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: testAccApigeeSecurityActionListQuery(orgId, envId),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_apigee_security_action.list_query", 0),
				},
			},
		},
	})
}

func testAccApigeeSecurityActionListQuery(orgId, envId string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_apigee_security_action" "list_query" {
  provider = google
  limit    = 1000
  config {
    org_id = %q
    env_id = %q
  }
}
`, orgId, envId)
}
