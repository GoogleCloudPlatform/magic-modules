package dns_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceDnsManagedZones_basic(t *testing.T) {
	// TODO: https://github.com/hashicorp/terraform-provider-google/issues/14158
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"name-1": fmt.Sprintf("tf-test-zone-%s", acctest.RandString(t, 10)),
		"name-2": fmt.Sprintf("tf-test-zone-%s", acctest.RandString(t, 10)),
	}

	project := envvar.GetTestProjectFromEnv()
	expectedId := fmt.Sprintf("projects/%s/managedZones", project)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducerFramework(t),
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccDataSourceDnsManagedZones_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dns_managed_zones.qa", "id", expectedId),
					resource.TestMatchResourceAttr("data.google_dns_managed_zones.qa", "managed_zones.#", regexp.MustCompile("^[1-9]")), // Non-zero number length
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZones_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dns_managed_zone" "one" {
  name        = "%{name-1}"
  dns_name    = "%{name-1}.hashicorptest.com."
  description = "tf test DNS zone"
}

resource "google_dns_managed_zone" "two" {
  name        = "%{name-2}"
  dns_name    = "%{name-2}.hashicorptest.com."
  description = "tf test DNS zone"
}

data "google_dns_managed_zones" "qa" {
}
`, context)
}
