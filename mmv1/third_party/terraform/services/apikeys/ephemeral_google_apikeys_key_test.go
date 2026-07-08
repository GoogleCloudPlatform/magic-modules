package apikeys_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/apikeys"
	_ "github.com/hashicorp/terraform-provider-google/google/services/secretmanager"
)

func TestAccEphemeralApikeysKey_basic(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	key := "tf-test-key-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralApikeysKey_setup(key),
			},
			{
				Config: testAccEphemeralApikeysKey_basic(key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.google_secret_manager_secret_version.from_full_name", "secret_data", "google_apikeys_key.key", "key_string"),
					resource.TestCheckResourceAttrPair("data.google_secret_manager_secret_version.from_key_id", "secret_data", "google_apikeys_key.key", "key_string"),
				),
			},
		},
	})
}

func testAccEphemeralApikeysKey_setup(key string) string {
	return fmt.Sprintf(`
resource "google_apikeys_key" "key" {
  name         = "%s"
  display_name = "Acceptance test API key"
}
`, key)
}

func testAccEphemeralApikeysKey_basic(key string) string {
	return fmt.Sprintf(`
resource "google_apikeys_key" "key" {
  name         = "%s"
  display_name = "Acceptance test API key"
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "%s"

  replication {
    auto {}
  }
}

ephemeral "google_apikeys_key" "from_full_name" {
  name = google_apikeys_key.key.id
}

resource "google_secret_manager_secret_version" "from_full_name" {
  secret                 = google_secret_manager_secret.secret.id
  secret_data_wo         = ephemeral.google_apikeys_key.from_full_name.key_string
  secret_data_wo_version = "1"
}

data "google_secret_manager_secret_version" "from_full_name" {
  secret  = google_secret_manager_secret_version.from_full_name.secret
  version = google_secret_manager_secret_version.from_full_name.version
}

ephemeral "google_apikeys_key" "from_key_id" {
  name    = google_apikeys_key.key.name
  project = google_apikeys_key.key.project
}

resource "google_secret_manager_secret_version" "from_key_id" {
  secret                 = google_secret_manager_secret.secret.id
  secret_data_wo         = ephemeral.google_apikeys_key.from_key_id.key_string
  secret_data_wo_version = "2"
}

data "google_secret_manager_secret_version" "from_key_id" {
  secret  = google_secret_manager_secret_version.from_key_id.secret
  version = google_secret_manager_secret_version.from_key_id.version
}
`, key, key)
}
