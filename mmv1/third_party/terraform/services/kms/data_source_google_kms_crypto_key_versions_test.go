package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleKmsCryptoKeyVersions_basic(t *testing.T) {
	asymSignKey := acctest.BootstrapKMSKeyWithPurpose(t, "ASYMMETRIC_SIGN")

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.versions", "versions.#", "2"),
			},
			// Asymmetric keys should have a public key
			{
				Config: testAccDataSourceGoogleKmsCryptoKeyVersions_basic(asymSignKey.CryptoKey.Name),
				Check:  resource.TestCheckResourceAttr("data.google_kms_crypto_key_versions.versions", "versions.0.public_key.0.state", "ENABLED"),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKeyVersions_basic(kmsKey string) string {
	return fmt.Sprintf(`
data "google_kms_crypto_key_versions" "versions" {
  crypto_key = "%s"
}
`, kmsKey)
}
