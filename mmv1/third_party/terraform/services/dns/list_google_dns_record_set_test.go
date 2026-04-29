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

func TestAccDNSRecordSetListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("tf-test-zone-%s", acctest.RandString(t, 10))
	recordLabel := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	project := envvar.GetTestProjectFromEnv()
	expectedName := fmt.Sprintf("%s.%s.hashicorptest.com.", recordLabel, zoneName)

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsRecordSetListResourceBasic(zoneName, recordLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_dns_record_set.foobar", "project", project),
					resource.TestCheckResourceAttr("google_dns_record_set.foobar", "managed_zone", zoneName),
					resource.TestCheckResourceAttr("google_dns_record_set.foobar", "name", expectedName),
					resource.TestCheckResourceAttr("google_dns_record_set.foobar", "type", "A"),
				),
			},
			{
				Query:  true,
				Config: testAccDnsRecordSetListQuery(zoneName),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_dns_record_set.all_in_zone", map[string]knownvalue.Check{
						"project":      knownvalue.StringExact(project),
						"managed_zone": knownvalue.StringExact(zoneName),
						"name":         knownvalue.StringExact(expectedName),
						"type":         knownvalue.StringExact("A"),
					}),
					querycheck.ExpectLengthAtLeast("google_dns_record_set.all_in_zone", 1),
				},
			},
		},
	})
}

func testAccDnsRecordSetListResourceBasic(zoneName, recordLabel string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "zone" {
  name     = %q
  dns_name = "%s.hashicorptest.com."
}

resource "google_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.zone.name
  name         = "%s.${google_dns_managed_zone.zone.dns_name}"
  type         = "A"
  ttl          = 300
  rrdatas      = ["192.168.1.0"]
}
`, zoneName, zoneName, recordLabel)
}

func testAccDnsRecordSetListQuery(zoneName string) string {
	return fmt.Sprintf(`
provider "google" {}

list "google_dns_record_set" "all_in_zone" {
  provider = google

  config {
    managed_zone = %q
  }
}
`, zoneName)
}
