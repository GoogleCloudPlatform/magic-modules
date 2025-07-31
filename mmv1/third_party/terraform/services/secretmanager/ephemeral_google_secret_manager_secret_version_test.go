// filepath: /Users/ramon/projects/personal/magic-modules/mmv1/third_party/terraform/services/secretmanager/ephemeral_google_secret_manager_secret_version_test.go
package secretmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccEphemeralSecretManagerSecretVersion_basic(t *testing.T) {
	t.Parallel()

	secret := "tf-test-secret-" + acctest.RandString(t, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralSecretManagerSecretVersion_basic(secret, "secret-data"),
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
`, secret, secretData)
}
