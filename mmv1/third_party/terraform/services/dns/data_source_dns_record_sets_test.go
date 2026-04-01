package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceDnsRecordSets_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsRecordSets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dns_record_sets.example", "rrsets.#", "1"),
					resource.TestCheckResourceAttr("data.google_dns_record_sets.example", "rrsets.0.type", "A"),
					resource.TestCheckResourceAttr("data.google_dns_record_sets.example", "rrsets.0.rrdatas.0", "192.168.1.0"),
				),
			},
		},
	})
}

func testAccDataSourceDnsRecordSets_basic(randString string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "zone" {
  name     = "tf-test-zone-%s"
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

data "google_dns_record_sets" "example" {
  managed_zone = google_dns_managed_zone.zone.name
  type         = "A"

  depends_on = [google_dns_record_set.rs]
}
`, randString, randString, randString)
}
