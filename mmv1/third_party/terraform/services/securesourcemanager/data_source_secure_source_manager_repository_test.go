package securesourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/securesourcemanager"
)

func TestAccSecureSourceManagerRepositoryDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecureSourceManagerRepositoryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecureSourceManagerRepositoryDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					// deletion_policy is a client-side virtual field not returned by the API; the test resource sets it to "DELETE" for cleanup.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_secure_source_manager_repository.default",
						"google_secure_source_manager_repository.default",
						[]string{"deletion_policy"},
					),
					// Verify omitted optional location and project fall back to provider defaults.
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_secure_source_manager_repository.without_optional_scope",
						"google_secure_source_manager_repository.default",
						[]string{"deletion_policy"},
					),
				),
			},
		},
	})
}

func testAccSecureSourceManagerRepositoryDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secure_source_manager_instance" "instance" {
  location        = "us-central1"
  instance_id     = "tf-test-my-instance%{random_suffix}"
  deletion_policy = "DELETE"
}

resource "google_secure_source_manager_repository" "default" {
  location        = "us-central1"
  repository_id   = "tf-test-my-repository%{random_suffix}"
  instance        = google_secure_source_manager_instance.instance.name
  description     = "Test Secure Source Manager Repository"
  deletion_policy = "DELETE"
}

data "google_secure_source_manager_repository" "default" {
  repository_id = google_secure_source_manager_repository.default.repository_id
  location      = google_secure_source_manager_repository.default.location
  project       = google_secure_source_manager_repository.default.project
}

data "google_secure_source_manager_repository" "without_optional_scope" {
  repository_id = google_secure_source_manager_repository.default.repository_id
}
`, context)
}
