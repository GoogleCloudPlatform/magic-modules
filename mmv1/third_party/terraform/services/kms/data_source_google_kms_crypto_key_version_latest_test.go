package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsCryptoKeyVersionLatest_basic(t *testing.T) {
	asymSignKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")
	asymDecrKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_DECRYPT")
	symKey := acctest.BootstrapKMSKey(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersionLatest_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version_latest.version_latest", "version", "2"),
			},
			// Asymmetric keys should have a public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersionLatest_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version_latest.version_latest", "public_key.#", "1"),
			},
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersionLatest_basic(asymDecrKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version_latest.version_latest", "public_key.#", "1"),
			},
			// Symmetric key should have no public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersionLatest_basic(symKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_version_latest.version_latest", "public_key.#", "0"),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKeyVersionLatest_basic(kmsKey string) string {
	return fmt.Sprintf(`
resource "google_kms_crypto_key_version" "version" {
	crypto_key = "%s"
	state = "ENABLED"
  }

data "google_kms_crypto_key_version_latest" "version_latest" {
  crypto_key = "%s"
}
`, kmsKey, kmsKey)
}
