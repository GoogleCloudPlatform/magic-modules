package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerDnsAuthorization_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerDnsAuthorizationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerDnsAuthorization_update0(context),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
			{
				Config: testAccCertificateManagerDnsAuthorization_update1(context),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerDnsAuthorization_update0(context map[string]interface{}) string {
	return acctest.Nprintf(`
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
	return acctest.Nprintf(`
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

func TestAccCertificateManagerDnsAuthorization_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestOrganizationTagKey(t, "certificate-manager-dns-auth-tagkey", map[string]interface{}{})

	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestOrganizationTagValue(t, "certificate-manager-dns-auth-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerDnsAuthorizationDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerDnsAuthorizationTags(context),
				Check: resource.TestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"google_certificate_manager_dns_authorization.default", "tags.%"),
				),
			},
			{
				ResourceName:            "google_certificate_manager_dns_authorization.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCertificateManagerDnsAuthorizationTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_dns_authorization" "default" {
        name          = "tf-test-dns-auth%{random_suffix}"
        description = "The default dns"
        labels = {
                a = "a"
        }
        domain          = "%{random_suffix}.hashicorptest.com"
	tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}
