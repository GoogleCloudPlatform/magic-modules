package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(t *testing.T) {
	asymSignKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")
	asymDecrKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_DECRYPT")
	symKey := acctest.BootstrapKMSKey(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.version", "version", "1"),
			},
			// Asymmetric keys should have a public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.version", "public_key.#", "1"),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(asymDecrKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.version", "public_key.#", "1"),
			},
			// Symmetric key should have no public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(symKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_latest_version.version", "public_key.#", "0"),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKeyLatestVersion_basic(kmsKey string) string {
	return fmt.Sprintf(`
data "google_kms_crypto_key_latest_version" "version" {
  crypto_key = "%s"
}
`, kmsKey)
}
