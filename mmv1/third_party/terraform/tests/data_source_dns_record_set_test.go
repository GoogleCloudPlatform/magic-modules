package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDnsRecordSet_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"google": func() (tfprotov5.ProviderServer, error) {
				provider, err := MuxedProviders(t.Name())
				return provider(), err
			},
		},
		CheckDestroy: testAccCheckDnsRecordSetDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsRecordSet_basic(randString(t, 10), randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_dns_record_set.rs", "google_dns_record_set.rs"),
				),
			},
		},
	})
}

func testAccDataSourceDnsRecordSet_basic(zoneName, recordSetName string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "zone" {
  name     = "test-zone"
  dns_name = "%s.hashicorptest.com."
}

resource "google_dns_record_set" "rs" {
  managed_zone = google_dns_managed_zone.zone.name
  name         = "%s.${google_dns_managed_zone.zone.dns_name}"
  type         = "A"
  ttl          = 300
  rrdatas      = [
	"192.168.1.0",
  ]
}
`, zoneName, recordSetName)
}

// Framework checkdestroy
func testAccCheckDnsRecordSetDestroyProducerFramework(t *testing.T) func(s *terraform.State) error {

	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dns_record_set" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			p := getTestFwProvider(t)

			url, err := replaceVarsForFrameworkTest(p, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if !p.ProdProvider.billingProject.IsNull() && p.ProdProvider.billingProject.String() != "" {
				billingProject = p.ProdProvider.billingProject.String()
			}

			_, diags := sendFrameworkRequest(&p.ProdProvider, "GET", billingProject, url, p.ProdProvider.userAgent, nil)
			if !diags.HasError() {
				return fmt.Errorf("DNSResourceDnsRecordSet still exists at %s", url)
			}
		}

		return nil
	}
}
