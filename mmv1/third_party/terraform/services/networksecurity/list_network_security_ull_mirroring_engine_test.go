package networksecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// TestAccNetworkSecurityUllMirroringEngineListQuery exercises the list resource
// for google_network_security_ull_mirroring_engine via a pure Query step.
// No infrastructure is created; it validates reachability of the list endpoint.
func TestAccNetworkSecurityUllMirroringEngineListQuery(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	project := envvar.GetTestProjectFromEnv()
	region := envvar.GetTestRegionFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: testAccNetworkSecurityUllMirroringEngineListQuery(project, region),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_network_security_ull_mirroring_engine.list_query", 0),
				},
			},
		},
	})
}

func testAccNetworkSecurityUllMirroringEngineListQuery(project, location string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_network_security_ull_mirroring_engine" "list_query" {
  provider = google
  limit    = 1000
  config {
    project  = %q
    location = %q
  }
}
`, project, location)
}
