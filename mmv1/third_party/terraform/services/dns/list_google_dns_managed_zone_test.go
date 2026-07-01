package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDnsManagedZoneListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	dnsName := fmt.Sprintf("tf-test-%s.hashicorptest.com.", acctest.RandString(t, 10))
	project := envvar.GetTestProjectFromEnv()
	t.Logf("Using project %s for testing", project)

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsManagedZoneListResource_basic(zoneName, dnsName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dns_managed_zone.foobar", "project", project),
					resource.TestCheckResourceAttr("google_dns_managed_zone.foobar", "name", zoneName),
					resource.TestCheckResourceAttr("google_dns_managed_zone.foobar", "dns_name", dnsName),
				),
			},
			{
				Query:  true,
				Config: testAccDnsManagedZoneListQuery(project),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_dns_managed_zone.all", map[string]knownvalue.Check{
						"name":    knownvalue.StringExact(zoneName),
						"project": knownvalue.StringExact(project),
					}),
					querycheck.ExpectLengthAtLeast("google_dns_managed_zone.all", 1),
				},
			},
		},
	})
}

func testAccDnsManagedZoneListResource_basic(name, dnsName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
  name       = %q
  dns_name   = %q
  visibility = "public"
}
`, name, dnsName)
}

func testAccDnsManagedZoneListQuery(project string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_dns_managed_zone" "all" {
  provider = google

  config {
    project = %q
  }
}
`, project)
}
