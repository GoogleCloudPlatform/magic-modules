package kms_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
)

func TestAccDataSourceGoogleKmsCryptoKey_basic(t *testing.T) {
	bootstrapped := kms.BootstrapKMSKey(t)

	// Name in the KMS client is in the format projects/<project>/locations/<location>/keyRings/<keyRingName>/cryptoKeys/<keyId>
	keyParts := strings.Split(bootstrapped.CryptoKey.Name, "/")
	cryptoKeyId := keyParts[len(keyParts)-1]

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsCryptoKey_basic(bootstrapped.KeyRing.Name, cryptoKeyId),
				Check:  resource.TestMatchResourceAttr("data.google_kms_crypto_key.kms_crypto_key", "id", regexp.MustCompile(bootstrapped.CryptoKey.Name)),
			},
		},
	})
}

func testAccDataSourceGoogleKmsCryptoKey_basic(keyRingName, cryptoKeyName string) string {
	return fmt.Sprintf(`
data "google_kms_crypto_key" "kms_crypto_key" {
  key_ring = "%s"
  name     = "%s"
}
`, keyRingName, cryptoKeyName)
}
