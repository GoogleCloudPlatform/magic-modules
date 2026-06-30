package kms_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/kms"
)

func TestAccEphemeralKMSKeyRing_basic(t *testing.T) {
	t.Parallel()

	keyRingName := "tf-test-keyring-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Step 1: create the key ring so it exists before the ephemeral resource opens.
				Config: testAccKMSKeyRing_basic(keyRingName),
			},
			{
				// Step 2: open the ephemeral resource against the already-existing key ring.
				Config: testAccEphemeralKMSKeyRing_basic(keyRingName),
			},
		},
	})
}

func testAccKMSKeyRing_basic(keyRingName string) string {
	return fmt.Sprintf(`
resource "google_kms_key_ring" "keyring" {
  name     = "%s"
  location = "global"
}
`, keyRingName)
}

func testAccEphemeralKMSKeyRing_basic(keyRingName string) string {
	return fmt.Sprintf(`
resource "google_kms_key_ring" "keyring" {
  name     = "%s"
  location = "global"
}

ephemeral "google_kms_key_ring" "ephemeral" {
  name     = google_kms_key_ring.keyring.name
  location = "global"
}
`, keyRingName)
}
