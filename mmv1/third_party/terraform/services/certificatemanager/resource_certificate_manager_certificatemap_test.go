package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerCertificateMap_tags(t *testing.T) {
	t.Parallel()
	tagKey := acctest.BootstrapSharedTestTagKey(t, "ccm-certificate-map-tagkey")
	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestTagValue(t, "ccm-certificate-map-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificateMapTags(context),
				Check: resource.TestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"google_certificate_manager_certificate_map.certificatemap", "tags.%"),
				),
			},
			{
				ResourceName:            "google_certificate_manager_certificate_map.certificatemap",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccCertificateManagerCertificateMapTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate_map" "certificatemap" {
  name = "tf-test-certificate-map-%{random_suffix}"
  description = "Global cert"
  tags = {
	"%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}
