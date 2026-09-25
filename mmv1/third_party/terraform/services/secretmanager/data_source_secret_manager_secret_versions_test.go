package secretmanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSecretManagerSecretVersions_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerSecretVersions_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_secret_manager_secret_versions.versions", "versions.#"),
					resource.TestCheckResourceAttrSet("data.google_secret_manager_secret_versions.versions", "versions.0.name"),
					resource.TestCheckResourceAttrSet("data.google_secret_manager_secret_versions.versions", "versions.0.version"),
					resource.TestCheckResourceAttr("data.google_secret_manager_secret_versions.versions", "versions.0.enabled", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSecretManagerSecretVersions_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret" {
  secret_id = "tf-test-secret-%{random_suffix}"
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "version" {
  secret = google_secret_manager_secret.secret.id
  secret_data = "my-secret-data"
}

data "google_secret_manager_secret_versions" "versions" {
  secret  = google_secret_manager_secret.secret.secret_id
  project = google_secret_manager_secret.secret.project
  depends_on = [google_secret_manager_secret_version.version]
}
`, context)
}
