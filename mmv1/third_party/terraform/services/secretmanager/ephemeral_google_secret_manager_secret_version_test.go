package secretmanager_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccEphemeralSecretManagerSecretVersion_basic(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	secret := "tf-test-secret-" + acctest.RandString(t, 10)
	secretData := "secret-data"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralSecretManagerSecretVersion_basic(secret, secretData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_secret_manager_secret_version.default", "secret_data", secretData),
				),
			},
		},
	})
}

func TestAccEphemeralSecretManagerSecretVersion_base64(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t)

	secret := "tf-test-secret-" + acctest.RandString(t, 10)
	secretData := "secret-data"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralSecretManagerSecretVersion_base64(secret, secretData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_secret_manager_secret_version.default", "is_secret_data_base64", "true"),
					resource.TestCheckResourceAttr("data.google_secret_manager_secret_version.default", "secret_data", base64.StdEncoding.EncodeToString([]byte(secretData))),
				),
			},
		},
	})
}

func testAccEphemeralSecretManagerSecretVersion_basic(secret, secretData string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "secret" {
  secret_id = "%s"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = "%s"
}

ephemeral "google_secret_manager_secret_version" "ephemeral" {
  secret  = google_secret_manager_secret_version.version.secret
  version = google_secret_manager_secret_version.version.version
}

resource "google_secret_manager_secret_version" "version_two_based_on_ephemeral" {
  secret  	             = google_secret_manager_secret_version.version.secret
  secret_data_wo 		 = ephemeral.google_secret_manager_secret_version.ephemeral.secret_data
  secret_data_wo_version = "1"
}

data "google_secret_manager_secret_version" "default" {
  secret  = google_secret_manager_secret_version.version_two_based_on_ephemeral.secret
  version = google_secret_manager_secret_version.version_two_based_on_ephemeral.version
}
`, secret, secretData)
}

func testAccEphemeralSecretManagerSecretVersion_base64(secret, secretData string) string {
	return fmt.Sprintf(`
resource "google_secret_manager_secret" "secret" {
  secret_id = "%s"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "version" {
  secret                 = google_secret_manager_secret.secret.id
  secret_data            = base64encode("%s")
  is_secret_data_base64  = true
}

ephemeral "google_secret_manager_secret_version" "ephemeral" {
  secret  				= google_secret_manager_secret_version.version.secret
  version 				= google_secret_manager_secret_version.version.version
  is_secret_data_base64 = true
}

resource "google_secret_manager_secret_version" "version_two_based_on_ephemeral" {
  secret  	             = google_secret_manager_secret_version.version.secret
  secret_data_wo 		 = ephemeral.google_secret_manager_secret_version.ephemeral.secret_data
  secret_data_wo_version = "1"
  is_secret_data_base64  = true
}

data "google_secret_manager_secret_version" "default" {
  secret                 = google_secret_manager_secret_version.version_two_based_on_ephemeral.secret
  version                = google_secret_manager_secret_version.version_two_based_on_ephemeral.version
  is_secret_data_base64  = true
}
`, secret, secretData)
}
