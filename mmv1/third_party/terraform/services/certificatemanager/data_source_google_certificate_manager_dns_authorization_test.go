package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCertificateManagerDnsAuthorizationDatasource(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerDnsAuthorizationDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_certificate_manager_dns_authorization.default", "google_certificate_manager_dns_authorization.default"),
				),
			},
		},
	})
}

func testAccCertificateManagerDnsAuthorizationDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_dns_authorization" "default" {
  name        = "tf-test-dns-auth-%{random_suffix}"
  location    = "global"
  description = "The default dns"
  domain      = "%{random_suffix}.hashicorptest.com"

}

data "google_certificate_manager_dns_authorization" "default" {
  name        = google_certificate_manager_dns_authorization.default.name
  domain      = "%{random_suffix}.hashicorptest.com"
  location    = "global"

}
`, context)
}
