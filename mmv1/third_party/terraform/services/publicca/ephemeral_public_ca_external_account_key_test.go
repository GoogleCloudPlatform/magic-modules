package publicca_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/publicca"
)

func TestAccEphemeralPublicCAExternalAccountKey_basic(t *testing.T) {
	t.Parallel()
	// Only boolean assertions are passed through the echo provider; key material
	// remains ephemeral. Omit project/location to exercise their default paths.
	config := `
ephemeral "google_public_ca_external_account_key" "test" {}
` + acctest.EchoResourceConfig(`{
  key_present = length(ephemeral.google_public_ca_external_account_key.test.key_id) > 0
  name_present = length(ephemeral.google_public_ca_external_account_key.test.name) > 0
  mac_decoded = base64decode(replace(replace(ephemeral.google_public_ca_external_account_key.test.b64_mac_key, "-", "+"), "_", "/")) == ephemeral.google_public_ca_external_account_key.test.mac_key
  url_encoding = trim(replace(replace(base64encode(ephemeral.google_public_ca_external_account_key.test.mac_key), "+", "-"), "/", "_"), "=") == ephemeral.google_public_ca_external_account_key.test.b64url_mac_key
}`)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(t),
		Steps: []resource.TestStep{{Config: config, Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.key_present", "true"),
			resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.name_present", "true"),
			resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.mac_decoded", "true"),
			resource.TestCheckResourceAttr(acctest.EchoResourceName, "data.url_encoding", "true"),
		)}},
	})
}
