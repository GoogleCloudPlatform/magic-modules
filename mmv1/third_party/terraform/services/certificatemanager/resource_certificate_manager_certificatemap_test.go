package certificatemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerCertificateMap_tags(t *testing.T) {
	t.Parallel()
	org := envvar.GetTestOrgFromEnv(t)
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	tagKey := acctest.BootstrapSharedTestTagKey(t, "ccm-certificates-tagkey")
	tagValue := acctest.BootstrapSharedTestTagValue(t, "ccm-certificates-tagvalue", tagKey)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccCertificateManagerCertificateTags(name, map[string]string{org + "/" + tagKey: tagValue}),
			},
			{
				ResourceName:            "google_certificate_manager_certificatemap.certificatemap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "tags"},
			},
		},
	})
}

func testAccCertificateManagerCertificateMapTags(name string, tags map[string]string) string {
	r := fmt.Sprintf(`
resource "google_certificate_manager_certificatemap" "certificatemap" {
  name = "tf-certificate-%s"
  description = "Global cert"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }
tags = {`, name)

	l := ""
	for key, value := range tags {
		l += fmt.Sprintf("%q = %q\n", key, value)
	}

	l += fmt.Sprintf("}\n}")
	return r + l
}

