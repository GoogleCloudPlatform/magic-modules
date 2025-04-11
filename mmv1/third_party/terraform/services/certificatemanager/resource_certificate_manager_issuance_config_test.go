package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccCertificateManagerIssuanceConfig_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestTagKey(t, "certificate_manager_issuance_config-tagkey")
	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestTagValue(t, "certificate_manager_issuance_config-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerIssuanceConfigTags(context),
			},
			{
				ResourceName:            "google_certificate_manager_certificate_issuance_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "tags"},
			},
		},
	})
}

func testAccCertificateManagerIssuanceConfigTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate_issuance_config" "default" {
  name        = "tf-test-issuance-config%{random_suffix}"
  description = "sample description for the issaunce config"
  location    = "us-central1"

  lifetime                    = "2592000s"
  key_algorithm               = "RSA_2048"
  rotation_window_percentage  = 80

  certificate_authority_config {
    certificate_authority_service_config {
      ca_pool = "projects/%{org}/locations/us-central1/caPools/tf-test-ca-pool%{random_suffix}"
    }
  }
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}
