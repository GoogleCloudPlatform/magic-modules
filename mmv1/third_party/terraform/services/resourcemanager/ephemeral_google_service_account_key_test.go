package resourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccEphemeralServiceAccountKey_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"project":       envvar.GetTestProjectFromEnv(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralServiceAccountKey_setup(context),
			},
			{
				Config: testAccEphemeralServiceAccountKey_basic(context),
			},
		},
	})
}

func testAccEphemeralServiceAccountKey_setup(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "test_account" {
  account_id   = "tf-test-%{random_suffix}"
  display_name = "Test Service Account"
}

resource "google_service_account_key" "key" {
  name = "${google_service_account.test_account.name}/keys/1234567890"
  public_key_type = "TYPE_RAW"
}
`, context)
}

func testAccEphemeralServiceAccountKey_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "test_account" {
  account_id   = "tf-test-%{random_suffix}"
  display_name = "Test Service Account"
}

resource "google_service_account_key" "key" {
  name = "${google_service_account.test_account.name}/keys/1234567890"
  public_key_type = "TYPE_RAW"
}

ephemeral "google_service_account_key" "key" {
  name            = "${google_service_account.test_account.name}/keys/1234567890"
  public_key_type = "TYPE_RAW"
}
`, context)
}
