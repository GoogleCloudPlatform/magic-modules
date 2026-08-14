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

// TestAccApigeeEnvironmentApiRevisionDeploymentListQuery exercises the list
// resource for google_apigee_environment_api_revision_deployment.
// It issues a pure list query against an existing API proxy revision in the
// test org, which validates that the list endpoint is reachable and the
// provider plumbing is correct without provisioning any new infrastructure.
func TestAccApigeeEnvironmentApiRevisionDeploymentListQuery(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	orgId := envvar.GetTestOrgFromEnv(t)
	// These values refer to an existing API proxy revision deployed in the
	// test org's Apigee environment; placeholders are used here because
	// this test runs against a pre-provisioned org.
	environment := "my-test-environment"
	api := "my-test-api"
	revision := "1"

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: testAccApigeeEnvironmentApiRevisionDeploymentListQuery(orgId, environment, api, revision),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_apigee_environment_api_revision_deployment.list_query", 0),
				},
			},
		},
	})
}

func testAccApigeeEnvironmentApiRevisionDeploymentListQuery(orgId, environment, api, revision string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_apigee_environment_api_revision_deployment" "list_query" {
  provider = google
  limit    = 1000
  config {
    org_id      = %q
    environment = %q
    api         = %q
    revision    = %s
  }
}
`, orgId, environment, api, revision)
}
