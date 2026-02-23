package discoveryengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us", "tftest-shared-key-"+acctest.RandString(t, 10)).CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_basic(context),
			},
			{
				ResourceName:            "google_discovery_engine_cmek_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cmek_config_id", "location", "project", "set_default"},
			},
		},
	})
}

func testAccDiscoveryEngineCmekConfig_discoveryengineCmekconfigDefaultExample_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_discovery_engine_cmek_config" "default" {
  location            = "us"
  cmek_config_id      = "tf-test-cmek-config-id%{random_suffix}"
  kms_key             = "%{kms_key_name}"
  set_default         = false
  depends_on = [google_kms_crypto_key_iam_member.crypto_key]
}

data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "crypto_key" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-discoveryengine.iam.gserviceaccount.com"
}
`, context)
}
