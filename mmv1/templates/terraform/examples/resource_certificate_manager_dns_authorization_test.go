package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCertificateManagerDnsAuthorization_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCertificateManagerDnsAuthorizationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerDnsAuthorization_update0(context),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: testAccCertificateManagerDnsAuthorization_update1(context),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func testAccCertificateManagerDnsAuthorization_update0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_certificate_manager_dns_authorization" "default" {
  name        = "tf-test-dns-auth%{random_suffix}"
  description = "The default dnss"
	labels = {
		a = "a"
	}
  domain      = "%{random_suffix}.hashicorptest.com"
}
`, context)
}

func testAccCertificateManagerDnsAuthorization_update1(context map[string]interface{}) string {
	return Nprintf(`
resource "google_certificate_manager_dns_authorization" "default" {
  name        = "tf-test-dns-auth%{random_suffix}"
  description = "The default dnss2"
	labels = {
		a = "b"
	}
  domain      = "%{random_suffix}.hashicorptest.com"
}
`, context)
}
