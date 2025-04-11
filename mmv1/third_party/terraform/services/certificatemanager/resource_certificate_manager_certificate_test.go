package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerCertificate_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestTagKey(t, "certificate_manager_certificate-tagkey")
	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestTagValue(t, "certificate_manager_certificate-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificateTags(context),
			},
			{
				ResourceName:            "google_certificate_manager_certificate.certificate",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccCertificateManagerCertificateTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate" "certificate" {
  name   = "tf-test-certificate-%{random_suffix}"
  description = "Global cert"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}
