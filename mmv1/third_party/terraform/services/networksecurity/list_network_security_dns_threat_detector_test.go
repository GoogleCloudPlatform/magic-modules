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

// TestAccNetworkSecurityDnsThreatDetectorListQuery exercises the list resource
// for google_network_security_dns_threat_detector via a pure Query step.
// No infrastructure is created; it validates reachability of the list endpoint.
func TestAccNetworkSecurityDnsThreatDetectorListQuery(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	project := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Query:  true,
				Config: testAccNetworkSecurityDnsThreatDetectorListQuery(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("google_network_security_dns_threat_detector.list_query", 0),
				},
			},
		},
	})
}

func testAccNetworkSecurityDnsThreatDetectorListQuery(project string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_network_security_dns_threat_detector" "list_query" {
  provider = google
  limit    = 1000
  config {
    project  = %q
    location = "global"
  }
}
`, project)
}
